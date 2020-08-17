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
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestOptionProfilingTime(t *testing.T) {

	tests := []struct {
		name                 string
		pr                   OptionProfilingTime
		wantErr              bool
		expectedQuerySeconds string
	}{
		{
			"setting to 30s",
			OptionProfilingTime(30 * time.Second),
			false,
			"30",
		},
		{
			"setting to 10 min",
			OptionProfilingTime(10 * time.Minute),
			false,
			"600",
		},
		{
			"setting to 0 seconds",
			OptionProfilingTime(0),
			false,
			"",
		},
		{
			"setting to 10 millseconds",
			OptionProfilingTime(10 * time.Millisecond),
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "localhost", nil)
			assert.NilError(t, err)
			if err := tt.pr.Apply(request); (err != nil) != tt.wantErr {
				t.Errorf("OptionProfilingTime.Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				query := request.URL.Query()
				val := query.Get("seconds")
				assert.Equal(t, tt.expectedQuerySeconds, val)
			}
		})
	}
}
