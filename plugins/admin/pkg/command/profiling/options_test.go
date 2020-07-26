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
