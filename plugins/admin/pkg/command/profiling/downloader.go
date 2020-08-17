// Copyright © 2020 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package profiling

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	pprofPort  uint32 = 8008
	localPort  uint32 = 18008
	secondsKey string = "seconds"
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

// Downloader struct holds all private fields
type Downloader struct {
	podName    string
	namespace  string
	readyCh    chan struct{} // closed by portforward.ForwardPorts() when connection is ready
	errorCh    chan error
	restConfig *rest.Config
	localPort  uint32
	client     *http.Client
	dialerFunc func(upgrader spdy.Upgrader, client *http.Client, method string, url *url.URL) httpstream.Dialer
}

// ProfileDownloader interface for profile downloader
type ProfileDownloader interface {
	Download(ProfileType, io.Writer, ...DownloadOptions) error
}

// RestConfigGetter interface to get restconfig
type RestConfigGetter interface {
	RestConfig() (*rest.Config, error)
}

// NewDownloader returns the profiling downloader and setup connections asynchronously
func NewDownloader(cfgGetter RestConfigGetter, podName, namespace string, endCh <-chan struct{}) (ProfileDownloader, error) {
	cfg, err := cfgGetter.RestConfig()
	if err != nil {
		return nil, err
	}
	d := &Downloader{
		podName:    podName,
		namespace:  namespace,
		readyCh:    make(chan struct{}),
		errorCh:    make(chan error),
		restConfig: cfg,
		localPort:  18008,
		client:     http.DefaultClient,
		dialerFunc: spdy.NewDialer,
	}
	err = d.connect(endCh)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// non-blocking call to forward remote port in pod to localhost
func (d *Downloader) connect(endCh <-chan struct{}) error {
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
	dialer := d.dialerFunc(upgrader, &http.Client{Transport: transport}, http.MethodPost, url)
	out := &bytes.Buffer{}
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", d.localPort, pprofPort)}, endCh, d.readyCh, out, out)
	if err != nil {
		return err
	}
	go func() {
		defer close(d.errorCh)
		// if the func ForwardPorts() returns, the connection should not be available.
		d.errorCh <- fw.ForwardPorts()
	}()
	return nil
}

// Download specific type of profile with options
func (d *Downloader) Download(t ProfileType, output io.Writer, options ...DownloadOptions) error {
	if t <= ProfileTypeUnknown || t >= ProfileType(len(ProfileEndpoints)) {
		return fmt.Errorf("unsupported profiling type %d", t)
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
			case <-d.errorCh:
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
		defer func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}()
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("download error: %s, code %d", string(body), resp.StatusCode)
		}
		_, err = io.Copy(output, resp.Body)
		if err != nil {
			return err
		}
		return nil
	case err := <-d.errorCh:
		return err
	}
}
