package ppcommon

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"time"
)

// Profile contains the metadata and the raw content of a pprof profile.
type Profile struct {
	// Email or ID of the user who owns this profile.
	User string
	// Name of the profile type - e.g: "heap".
	Name string
	// Content is the raw profile data content, in pprof format. It's a slice
	// of bytes instead of a string or []string to avoid the 500-bytes limit
	// that StringProperty have. Reference:
	// http://stackoverflow.com/questions/11178869/overcome-appengine-500-byte-string-limit-in-python-consider-text
	Content []byte
	// Time when the profile was collected.
	Time time.Time
}

func RawProfile(p string) (raw string, err error) {
	buf := new(bytes.Buffer)
	funcs := map[string]string{}
	re := regexp.MustCompile(`#\s+(0x[^ ]+)\s+([^+ ]+)[+ ]`)

	fmt.Fprintf(buf, "--- symbol\nbinary=unknown\n")

	scanner := bufio.NewScanner(bytes.NewBufferString(p))

	innerBuf := new(bytes.Buffer)
	for scanner.Scan() {
		line := scanner.Text()
		m := re.FindStringSubmatch(line)
		if m != nil {
			funcs[m[1]] = m[2]
		} else {
			fmt.Fprintf(innerBuf, "%v\n", line)
		}
	}

	for pc, f := range funcs {
		fmt.Fprintf(buf, "%v %v\n", pc, f)
	}
	fmt.Fprintf(buf, "---\n%v", innerBuf)
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return buf.String(), nil

}
