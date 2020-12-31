package version

import (
	"time"
)

var (
	version = "0.0.0"
	commit  = "HEAD"
	date    = "1970-01-01T00:00:00Z"
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Date() time.Time {
	d, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Unix(0, 0)
	}

	return d
}
