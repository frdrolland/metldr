// tqa_pcap_tool project main.go
package main

import (
	"fmt"
	"regexp"
	//	"time"

	"github.com/frdrolland/pcaptool/cfg"
	"github.com/frdrolland/pcaptool/cli"
	"github.com/frdrolland/pcaptool/parsing"
)

// embed regexp.Regexp in a new type so we can extend it
type extRegexp struct {
	*regexp.Regexp
}

var (
	config cfg.Configuration
)

// add a new method to our new regular expression type
func (r *extRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}

// Verbose prints on standard output string argument if (and only if) -v option has been set in program arguments.
func Verbose(s string, args ...interface{}) {
	if config.Verbose {
		fmt.Printf("%s\n", s)
	}
}

func main() {
	//	t := time.Now()
	//t := time.Unix(0, 1501681526043505000)

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	Verbose("Parsed arguments : %s\n", config)
	//	Verbose("config = %s\n", config)

	// Parse log and extract JSON from each line for specific Regexp
	for _, currFile := range config.Files {
		// element is the element from someSlice for where we are
		Verbose("Importing file %s\n", currFile)
		parsing.ParseLines(currFile, func(s string) (string, bool) {

			// ok:
			re := regexp.MustCompile("(?P<timestamp>.*)\\s+\\|\\s+(.*)\\s+\\|\\s+(.*)\\s+\\|\\s+(.*)\\s+\\|\\s+(.*)\\s+\\|\\s+Connectors\\.hpp\\:\\d+\\s+\\|\\s+(?P<json>(.*))")

			n1 := re.SubexpNames()
			r2 := re.FindStringSubmatch(s)
			if nil == r2 {
				return "", false
			}

			md := map[string]string{}

			for i, n := range r2 {
				md[n1[i]] = n
			}
			//fmt.Printf("[debug] timestamp = %s\n", md["timestamp"])
			//fmt.Printf("md = %v\n", md)
			json := md["json"]

			if "" == json {
				return "", false
			}
			return json, true
		})
	}
}
