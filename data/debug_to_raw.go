// Converts a debug heap profile to a "raw" profile.
//
// The debug profile produced by Go has the function names next to their
// addresses as inline comments in the profile. We use those to generate a
// "symbols" section and create a valid "raw" profile. This appears to work for
// heap profiles at least.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

func main() {

	funcs := map[string]string{}
	buf := new(bytes.Buffer)

	fmt.Println("--- symbol\nbinary=unknown\n")
	scanner := bufio.NewScanner(os.Stdin)
	re := regexp.MustCompile(`#\s+(0x[^ ]+)\s+([^+ ]+)[+ ]`)

	for scanner.Scan() {
		line := scanner.Text()
		m := re.FindStringSubmatch(line)
		if m != nil {
			funcs[m[1]] = m[2]
		} else {
			fmt.Fprintf(buf, "%v\n", line)
		}
	}
	for pc, f := range funcs {
		fmt.Printf("%v %v\n", pc, f)
	}
	fmt.Printf("---\n")
	fmt.Print(buf)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
