package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

// DO NOT COMMIT
// sample code to download profile using profile downloader
func main() {

	kubeconfig := os.Getenv("KUBECONFIG")
	// If we have an explicit indication of where the kubernetes config lives, read that.
	if kubeconfig == "" {
		log.Fatalf("please set KUBECONFIG")
	}
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	// set signal handler and to close this channel to cancel all requests
	end := make(chan struct{})
	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// find autoscaler pod name
	pods, err := c.CoreV1().Pods("knative-serving").List(metav1.ListOptions{
		LabelSelector: "app=autoscaler",
	})
	if err != nil || len(pods.Items) == 0 {
		log.Fatal(err)
	}

	podName := pods.Items[0].Name
	d, err := NewDownloader(cfg, podName, "knative-serving", end)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("getting profile for pod %s", podName)
	if err != nil {
		log.Fatal(err)
	}
	f1, _ := ioutil.TempFile("", "")
	defer f1.Close()
	log.Println("downloading 5s CPU profile")
	err = d.Download(ProfileTypeProfile, f1, ProfilingTime(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cpu profile saved to %s", f1.Name())
	f2, _ := ioutil.TempFile("", "")
	defer f2.Close()
	log.Println("downloading heap profile")
	err = d.Download(ProfileTypeHeap, f2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("heap saved to %s", f2.Name())
}
