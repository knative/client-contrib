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
	"net/http"
	"strconv"
	"time"
)

// DownloadOptions interface to manipulate the http request to pprof server
type DownloadOptions interface {
	Apply(*http.Request) error
}

// OptionProfilingTime is option to add a seconds param in the http request
type OptionProfilingTime time.Duration

var _ DownloadOptions = OptionProfilingTime(time.Second)

// Apply implements DownloadOptions interface for type ProfilingTime
func (pr OptionProfilingTime) Apply(req *http.Request) error {
	query := req.URL.Query()
	seconds := int64(time.Duration(pr) / time.Second)
	if seconds <= 0 {
		return nil
	}
	query.Set(secondsKey, strconv.FormatInt(seconds, 10))
	req.URL.RawQuery = query.Encode()
	return nil
}
