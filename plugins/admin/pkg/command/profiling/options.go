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
