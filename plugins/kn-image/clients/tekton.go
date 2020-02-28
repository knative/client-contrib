package clients

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tekton_v1alpha1_client "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"time"
)

const (
	BuildpacksBuilderName  = "buildpacks-v3"
	BuildpacksBuilderImage = "cloudfoundry/cnb:bionic"
	KanikoBuilderName      = "kaniko"
	KanikoBuilderImage     = "gcr.io/kaniko-project/executor:v0.13.0"
	BuildTimeout           = 300
)

type TektonClient interface {
	// Create git resource struct from provided options
	ConstructGitResource(name, url, revision, namespace string) *pipelinev1alpha1.PipelineResource

	// Create image resource struct from provided options
	ConstructImageResource(name, image, builder, serviceAccount, namespace string) *pipelinev1alpha1.PipelineResource

	// Get a PipelineResource by its unique name
	GetPipelineResource(name string) (*v1alpha1.PipelineResource, error)

	// Check if the PipelineResource exists
	PipelineResourceExists(name string) (bool, error)

	// Create a new pipelineresource
	CreatePipelineResource(pipelineresource *v1alpha1.PipelineResource) error

	// Update the given pipelineresource
	UpdatePipelineResource(pipelineresource *v1alpha1.PipelineResource) error

	// Check if the builder task exists
	TaskExists(name string) error

	// Create task run struct from provided options
	ConstructTaskRun(name, builder, gitResourceName, imageResourceName, serviceAccount, namespace string) *pipelinev1alpha1.TaskRun

	// Get a task run by its unique name
	GetTaskRun(name string) (*v1alpha1.TaskRun, error)

	// Create a new task run
	StartTaskRun(task *v1alpha1.TaskRun) (string, error)

	// Build image from git
	BuildFromGit(name, builder, gitUrl, gitRevision, image, serviceAccount, namespace string) error
}

type tektonClient struct {
	client    tekton_v1alpha1_client.TektonV1alpha1Interface
	namespace string
}

// Create a new client facade for the provided namespace
func NewTektonClient(client tekton_v1alpha1_client.TektonV1alpha1Interface, namespace string) TektonClient {
	return &tektonClient{
		client:    client,
		namespace: namespace,
	}
}

func (cl *tektonClient) TaskExists(name string) error {
	_, err := cl.client.Tasks(cl.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

// ArrayOrString creates an ArrayOrString of type ParamTypeString or ParamTypeArray, based on
// how many inputs are given (>1 input will create an array, not string).
func ArrayOrString(value string, additionalValues ...string) *pipelinev1alpha1.ArrayOrString {
	if len(additionalValues) > 0 {
		additionalValues = append([]string{value}, additionalValues...)
		return &pipelinev1alpha1.ArrayOrString{
			Type:     pipelinev1alpha1.ParamTypeArray,
			ArrayVal: additionalValues,
		}
	}
	return &pipelinev1alpha1.ArrayOrString{
		Type:      pipelinev1alpha1.ParamTypeString,
		StringVal: value,
	}
}

// Create git resource struct from provided options
func (cl *tektonClient) ConstructGitResource(name, url, revision, namespace string) *pipelinev1alpha1.PipelineResource {

	gitresource := pipelinev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: pipelinev1alpha1.PipelineResourceSpec{
			Type: pipelinev1alpha1.PipelineResourceTypeGit,
			Params: []pipelinev1alpha1.ResourceParam{{
				Name:  "url",
				Value: url,
			}, {
				Name:  "revision",
				Value: revision,
			}},
		},
	}

	return &gitresource
}

// Create image resource struct from provided options
func (cl *tektonClient) ConstructImageResource(name, image, builder, serviceAccount, namespace string) *pipelinev1alpha1.PipelineResource {

	imageresource := pipelinev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				// The builder and service account will be stored in image output resource labels
				"builder": builder,
				"serviceaccount": serviceAccount,
			},
		},
		Spec: pipelinev1alpha1.PipelineResourceSpec{
			Type: pipelinev1alpha1.PipelineResourceTypeImage,
			Params: []pipelinev1alpha1.ResourceParam{{
				Name:  "url",
				Value: image,
			}},
		},
	}

	return &imageresource
}

func (cl *tektonClient) PipelineResourceExists(name string) (bool, error) {
	_, err := cl.client.PipelineResources(cl.namespace).Get(name, metav1.GetOptions{})
	if api_errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Get a pipelineresource by its unique name
func (cl *tektonClient) GetPipelineResource(name string) (*v1alpha1.PipelineResource, error) {
	pipelineresource, err := cl.client.PipelineResources(cl.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pipelineresource, nil
}

// Create a new pipelineresource
func (cl *tektonClient) CreatePipelineResource(pipelineresource *v1alpha1.PipelineResource) error {
	_, err := cl.client.PipelineResources(cl.namespace).Create(pipelineresource)
	if err != nil {
		return err
	}
	return nil
}

// Update the given pipelineresource
func (cl *tektonClient) UpdatePipelineResource(pipelineresource *v1alpha1.PipelineResource) error {
	var retries = 0
	var MaxUpdateRetries = 3
	for {
		existingPipelineResource, err := cl.client.PipelineResources(cl.namespace).Get(pipelineresource.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updateResource := pipelineresource.DeepCopy()
		updateResource.ResourceVersion = existingPipelineResource.ResourceVersion

		_, err = cl.client.PipelineResources(cl.namespace).Update(updateResource)
		if err != nil {
			// Retry to update when a resource version conflict exists
			if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
				retries++
				continue
			}
			return err
		}
		return nil
	}
}

// Create Task struct from provided options
func (cl *tektonClient) ConstructTaskRun(name, builder, gitResourceName, imageResourceName, serviceAccount, namespace string) *pipelinev1alpha1.TaskRun {

	builderImage := ""
	if builder == BuildpacksBuilderName {
		builderImage = BuildpacksBuilderImage
	} else if builder == KanikoBuilderName {
		builderImage = KanikoBuilderImage
	}

	return &pipelinev1alpha1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: name + "-build-",
		},
		Spec: pipelinev1alpha1.TaskRunSpec{
			TaskRef: &pipelinev1alpha1.TaskRef{
				Name: builder,
			},
			Inputs: pipelinev1alpha1.TaskRunInputs{
				Resources: []pipelinev1alpha1.TaskResourceBinding{{
					ResourceRef: pipelinev1alpha1.PipelineResourceRef{
						Name:       gitResourceName,
					},
					Name: "source",
				}},
				Params: []pipelinev1alpha1.Param{
					{
						Name:  "BUILDER_IMAGE",
						Value: *ArrayOrString(builderImage),
					},
				},
			},
			Outputs: pipelinev1alpha1.TaskRunOutputs{
				Resources: []pipelinev1alpha1.TaskResourceBinding{{
					ResourceRef: pipelinev1alpha1.PipelineResourceRef{
						Name:       imageResourceName,
					},
					Name: "image",
				}},
			},
			ServiceAccount: serviceAccount,
		},
	}
	return nil
}

// Get a tas krun by its unique name
func (cl *tektonClient) GetTaskRun(name string) (*v1alpha1.TaskRun, error) {
	taskRun, err := cl.client.TaskRuns(cl.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return taskRun, nil
}

// Start a new task run
func (cl *tektonClient) StartTaskRun(taskrun *v1alpha1.TaskRun) (string, error) {
	newTaskRun, err := cl.client.TaskRuns(cl.namespace).Create(taskrun)
	if err != nil {
		return "", err
	}

	time.Sleep(5 * time.Second)
	i := 0
	for  i < BuildTimeout {
		taskrun, err = cl.client.TaskRuns(cl.namespace).Get(newTaskRun.Name, metav1.GetOptions{})
		if taskrun.Status.Conditions[0].Type == apis.ConditionSucceeded && taskrun.Status.Conditions[0].Status == "True" {
			fmt.Println("[INFO] Build task run", taskrun.Name, "is ready from", taskrun.Status.StartTime, "to", taskrun.Status.CompletionTime)
			return taskrun.Name, nil
		} else {
			fmt.Println("[INFO] Build task run", taskrun.Name, "is still", taskrun.Status.Conditions[0].Reason, ", waiting")
			time.Sleep(5 * time.Second)
		}
		i += 5
		time.Sleep(5 * time.Second)
	}

	return taskrun.Name, fmt.Errorf("the build taskrun is not ready after timeout")
}


// Build application from git
func (cl *tektonClient) BuildFromGit(name, builder, gitUrl, gitRevision, image, serviceAccount, namespace string) error {

	fmt.Println("[INFO] Building image", name, "in namespace", namespace)
	fmt.Println("[INFO] By using builder", builder, "and service account", serviceAccount)
	fmt.Println("[INFO] From git repo", gitUrl, "revision", gitRevision)

	err := cl.TaskExists(builder)
	if err != nil {
		return err
	}
	fmt.Println("[INFO] Get task", builder, "successfully")

	// Create git resource
	gitResourceName := name + "-git"
	gitResource := cl.ConstructGitResource(gitResourceName, gitUrl, gitRevision, namespace)
	gitResourceExists, err := cl.PipelineResourceExists(gitResourceName)
	if err != nil {
		glog.Fatalf("[ERROR] Check git resource exist error: %s", err)
	}
	if gitResourceExists {
		fmt.Println("[INFO] git resource", gitResourceName, "exists, updating")
		err = cl.UpdatePipelineResource(gitResource)
	} else {
		fmt.Println("[INFO] git resource", gitResourceName, "doesn't exist, creating")
		err = cl.CreatePipelineResource(gitResource)
	}
	if err != nil {
		glog.Fatalf("[ERROR] Create git resource error: %s", err)
	}

	// Create image resource
	imageResourceName := name + "-image"
	imageResource := cl.ConstructImageResource(imageResourceName, image, builder, serviceAccount, namespace)
	imageResourceExists, err := cl.PipelineResourceExists(imageResourceName)
	if err != nil {
		glog.Fatalf("[ERROR] Check image resource exist error: %s", err)
	}
	if imageResourceExists {
		fmt.Println("[INFO] Image resource", imageResourceName, "exists, updating")
		err = cl.UpdatePipelineResource(imageResource)
	} else {
		fmt.Println("[INFO] Image resource", imageResourceName, "doesn't exist, creating")
		err = cl.CreatePipelineResource(imageResource)
	}
	if err != nil {
		glog.Fatalf("[ERROR] Create image resource error: %s", err)
	}

	// Start source-to-image task run
	buildTaskRun := cl.ConstructTaskRun(name, builder, gitResourceName, imageResourceName, serviceAccount, namespace)
	fmt.Println("[INFO] Start task run for image", name)
	taskRunName, err := cl.StartTaskRun(buildTaskRun)
	if err != nil {
		return err
	}
	fmt.Println("[INFO] Complete task run", taskRunName,"successfully")
	fmt.Println("[INFO] Complete building application", name, "image in namespace", namespace)

	return nil
}
