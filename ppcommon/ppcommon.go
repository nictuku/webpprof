package ppcommon

import (
	"time"
)

// profile contains the metadata and the raw content of a pprof profile.
type Profile struct {
	// Name of the profile type - e.g: "heap".
	Name string
	// Content is the raw profile data content, in pprof format. It's a slice
	// of bytes instead of a string or []string to avoid the 500-bytes limit
	// that StringProperty have. Reference:
	// http://stackoverflow.com/questions/11178869/overcome-appengine-500-byte-string-limit-in-python-consider-text
	Content []byte
	Time    time.Time
}
