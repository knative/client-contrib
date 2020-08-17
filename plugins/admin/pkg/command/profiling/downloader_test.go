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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
)

func fakeDialerFunc(t *testing.T, expectedURL *url.URL, exceptedError error) func(upgrader spdy.Upgrader, client *http.Client, method string, url *url.URL) httpstream.Dialer {
	return func(upgrader spdy.Upgrader, client *http.Client, method string, url *url.URL) httpstream.Dialer {
		return &fakeDialer{
			t:             t,
			url:           url,
			exceptedURL:   expectedURL,
			exceptedError: exceptedError,
		}
	}
}

type fakeDialer struct {
	t   *testing.T
	url *url.URL

	exceptedURL   *url.URL
	exceptedError error
}

type fakeConnection struct {
	closed    bool
	closeChan chan bool
}

func (f *fakeDialer) Dial(protocols ...string) (httpstream.Connection, string, error) {
	if f.exceptedURL != nil {
		assert.DeepEqual(f.t, f.exceptedURL, f.url)
	}
	if f.exceptedError != nil {
		return nil, "", f.exceptedError
	}
	return &fakeConnection{
		closed:    false,
		closeChan: make(chan bool),
	}, "", nil
}

func (f *fakeConnection) CreateStream(headers http.Header) (httpstream.Stream, error) {
	return nil, nil
}

func (f *fakeConnection) Close() error {
	if !f.closed {
		f.closed = true
		close(f.closeChan)
	}
	return nil
}

func (f *fakeConnection) CloseChan() <-chan bool {
	return f.closeChan

}
func (f *fakeConnection) SetIdleTimeout(timeout time.Duration) {
	// no-op
}

func TestProfileDownloader_connect(t *testing.T) {
	t.Run("connect success", func(t *testing.T) {
		d := &Downloader{
			podName:   "pod-1",
			namespace: "mynamespace",
			readyCh:   make(chan struct{}),
			errorCh:   make(chan error),
			client:    http.DefaultClient,
			localPort: 12345,
			restConfig: &rest.Config{
				Host: "http://localhost:12345",
			},
			dialerFunc: fakeDialerFunc(t,
				&url.URL{
					Scheme: "http",
					Host:   "localhost:12345",
					Path:   "/api/v1/namespaces/mynamespace/pods/pod-1/portforward",
				},
				nil),
		}
		ch := make(chan struct{})
		err := d.connect(ch)
		assert.NilError(t, err)
		// should be ready
		<-d.readyCh
		close(ch)
		err = <-d.errorCh
		assert.NilError(t, err)
	})

	t.Run("dial error", func(t *testing.T) {
		exceptDialError := errors.New("dial error")
		d := &Downloader{
			podName:   "pod-1",
			namespace: "mynamespace",
			readyCh:   make(chan struct{}),
			errorCh:   make(chan error),
			client:    http.DefaultClient,
			localPort: 12345,
			restConfig: &rest.Config{
				Host: "http://localhost:12345",
			},
			dialerFunc: fakeDialerFunc(t,
				&url.URL{
					Scheme: "http",
					Host:   "localhost:12345",
					Path:   "/api/v1/namespaces/mynamespace/pods/pod-1/portforward",
				},
				exceptDialError),
		}
		ch := make(chan struct{})
		err := d.connect(ch)
		assert.NilError(t, err)
		go func() {
			select {
			case <-d.readyCh:
				t.Error("ready chan should not be closed")
			case <-ch:
			}
		}()
		// should receive error before we close ch
		err = <-d.errorCh
		assert.ErrorContains(t, err, exceptDialError.Error())
		close(ch)
	})
}

func TestProfileDownload(t *testing.T) {
	t.Run("download heap profile success", func(t *testing.T) {
		downloadData := []byte("some-binary-data")
		server := httptest.NewServer(http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, "/debug/pprof/heap", req.URL.RequestURI())
				rw.WriteHeader(http.StatusOK)
				rw.Write(downloadData)
			},
		))
		defer server.Close()

		listenerAddr := server.Listener.Addr().String()
		_, portString, err := net.SplitHostPort(listenerAddr)
		assert.NilError(t, err)
		port, err := strconv.ParseInt(portString, 10, 0)
		assert.NilError(t, err)

		d := &Downloader{
			readyCh:   make(chan struct{}),
			errorCh:   make(chan error),
			client:    http.DefaultClient,
			localPort: uint32(port),
		}
		errChan := make(chan error)
		output := &bytes.Buffer{}
		go func() {
			errChan <- d.Download(ProfileTypeHeap, output)
		}()
		close(d.readyCh)

		err = <-errChan
		assert.NilError(t, err)

		bs, err := ioutil.ReadAll(output)
		assert.NilError(t, err)
		assert.DeepEqual(t, downloadData, bs)
	})

	t.Run("download error caused by response code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
				io.WriteString(rw, "not found")
			},
		))
		defer server.Close()

		listenerAddr := server.Listener.Addr().String()
		_, portString, err := net.SplitHostPort(listenerAddr)
		assert.NilError(t, err)
		port, err := strconv.ParseInt(portString, 10, 0)
		assert.NilError(t, err)

		d := &Downloader{
			readyCh:   make(chan struct{}),
			errorCh:   make(chan error),
			client:    http.DefaultClient,
			localPort: uint32(port),
		}
		errChan := make(chan error)
		output := &bytes.Buffer{}
		go func() {
			errChan <- d.Download(ProfileTypeHeap, output)
		}()
		close(d.readyCh)

		err = <-errChan
		assert.ErrorContains(t, err, "download error: not found, code 404")
	})

	t.Run("unsupported profile type", func(t *testing.T) {

		d := &Downloader{
			readyCh: make(chan struct{}),
			errorCh: make(chan error),
			client:  http.DefaultClient,
		}
		errChan := make(chan error)
		output := &bytes.Buffer{}

		go func() {
			errChan <- d.Download(ProfileTypeUnknown, output)
		}()
		close(d.readyCh)
		var err error
		err = <-errChan
		assert.ErrorContains(t, err, "unsupported profiling type")

		go func() {
			errChan <- d.Download(ProfileType(len(ProfileEndpoints)), output)
		}()
		err = <-errChan
		assert.ErrorContains(t, err, "unsupported profiling type")
	})

	t.Run("error occoured while download is not started", func(t *testing.T) {
		d := &Downloader{
			readyCh: make(chan struct{}),
			errorCh: make(chan error),
			client:  http.DefaultClient,
		}
		errChan := make(chan error)
		output := &bytes.Buffer{}

		go func() {
			errChan <- d.Download(ProfileTypeHeap, output)
		}()

		e := fmt.Errorf("dummy connection error")
		d.errorCh <- e
		var err error
		err = <-errChan
		assert.Error(t, err, e.Error())
	})

	t.Run("request canceled while download is started", func(t *testing.T) {
		downloadData := []byte("some-binary-data")
		server := httptest.NewServer(http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusOK)
				for _, b := range downloadData {
					<-time.After(time.Second) // write at 1 byte/second
					rw.Write([]byte{b})
				}
			},
		))
		defer server.Close()

		listenerAddr := server.Listener.Addr().String()
		_, portString, err := net.SplitHostPort(listenerAddr)
		assert.NilError(t, err)
		port, err := strconv.ParseInt(portString, 10, 0)
		assert.NilError(t, err)

		d := &Downloader{
			readyCh:   make(chan struct{}),
			errorCh:   make(chan error),
			client:    http.DefaultClient,
			localPort: uint32(port),
		}
		errChan := make(chan error)
		output := &bytes.Buffer{}
		go func() {
			errChan <- d.Download(ProfileTypeHeap, output)
		}()
		close(d.readyCh)
		<-time.After(1 * time.Second)
		close(d.errorCh)

		err = <-errChan
		assert.ErrorContains(t, err, "context canceled")
	})
}
