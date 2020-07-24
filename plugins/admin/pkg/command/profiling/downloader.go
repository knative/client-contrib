package profiling

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	pprofPort  uint32 = 8008
	secondsKey        = "seconds"
)

// ProfileType enums for all supported profiles
type ProfileType int

const (
	ProfileTypeUnknown ProfileType = iota
	ProfileTypeHeap
	ProfileTypeProfile
	ProfileTypeBlock
	ProfileTypeTrace
	ProfileTypeAllocs
	ProfileTypeMutex
	ProfileTypeGoroutine
	ProfileTypeThreadCreate
)

// ProfileEndpoints array maps ProfileType to the string endpoint for pprof
var ProfileEndpoints = [...]string{
	ProfileTypeUnknown:      "",
	ProfileTypeHeap:         "heap",
	ProfileTypeProfile:      "profile",
	ProfileTypeBlock:        "block",
	ProfileTypeTrace:        "trace",
	ProfileTypeAllocs:       "allocs",
	ProfileTypeMutex:        "mutex",
	ProfileTypeGoroutine:    "goroutine",
	ProfileTypeThreadCreate: "threadcreate",
}

// DownloadOptions interface to manipulate the http request to pprof server
type DownloadOptions interface {
	Apply(*http.Request) error
}

// Downloader struct holds all private fields
type Downloader struct {
	podName    string
	namespace  string
	stopCh     <-chan struct{} // Close this will trigger closing for the endCh create by our own and then cancel all sub goroutines
	endCh      chan struct{}
	readyCh    chan struct{}
	restConfig *rest.Config
	localPort  uint32
	client     *http.Client
}

// NewDownloader returns the profiling downloader and setup connections asynchronously
func NewDownloader(cfg *rest.Config, podName, namespace string, endCh <-chan struct{}) (*Downloader, error) {
	d := &Downloader{
		podName:    podName,
		namespace:  namespace,
		readyCh:    make(chan struct{}),
		stopCh:     endCh,
		restConfig: cfg,
		localPort:  18008,
		client:     http.DefaultClient,
	}
	err := d.connect()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// non-blocking call to forward remote port in pod to localhost
func (d *Downloader) connect() error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", d.namespace, d.podName)
	transport, upgrader, err := spdy.RoundTripperFor(d.restConfig)
	if err != nil {
		return err
	}
	u, err := url.Parse(d.restConfig.Host)
	if err != nil {
		return err
	}
	url := &url.URL{
		Host:   u.Host,
		Scheme: u.Scheme,
		Path:   path,
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, url)
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", d.localPort, pprofPort)}, d.stopCh, d.readyCh, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}
	go func() {
		defer close(d.endCh)
		err := fw.ForwardPorts()
		// if the func ForwardPorts() returns, the connection should not be available.
		if err != nil {
			// TODO: Log for error?
		}
	}()
	return nil
}

// ProfilingTime is option to add a seconds param in the http request
type ProfilingTime time.Duration

// Apply implements DownloadOptions interface for type ProfilingTime
func (pr ProfilingTime) Apply(req *http.Request) error {
	query := req.URL.Query()
	seconds := int64(time.Duration(pr) / time.Second)
	query.Set(secondsKey, strconv.FormatInt(seconds, 10))
	req.URL.RawQuery = query.Encode()
	return nil
}

// Download specific type of profile with options
func (d *Downloader) Download(t ProfileType, output io.Writer, options ...DownloadOptions) error {
	if t <= ProfileTypeUnknown || t >= ProfileType(len(ProfileEndpoints)) {
		return fmt.Errorf("unknown profiling type %d", t)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case <-d.readyCh:
		// connection ready
		var err error
		url := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("127.0.0.1:%d", d.localPort),
			Path:   fmt.Sprintf("/debug/pprof/%s", ProfileEndpoints[t]),
		}
		go func() {
			select {
			// request succeeded
			case <-ctx.Done():
				break
			case <-d.endCh:
				cancel()
			}
		}()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
		if err != nil {
			return err
		}
		for _, o := range options {
			if err = o.Apply(req); err != nil {
				return err
			}
		}
		resp, err := d.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		_, err = io.Copy(output, resp.Body)
		if err != nil {
			return err
		}
		return nil
	case <-d.endCh:
		return nil
	}
}
