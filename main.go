// tqa_pcap_tool project main.go
package main

import (
	"fmt"
	"regexp"

	"github.com/frdrolland/pcaptool/cli"
	"github.com/frdrolland/pcaptool/parsing"
)

// embed regexp.Regexp in a new type so we can extend it
type extRegexp struct {
	*regexp.Regexp
}

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

func main() {

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	fmt.Printf("config = %s\n", config)

	// Parse log
	for _, currFile := range config.Files {
		// element is the element from someSlice for where we are
		parsing.ParseLines(currFile, func(s string) (string, bool) {

			//re := regexp.MustCompile("(.*)\\s\\|\\s(.*)\\s\\|\\s(.*)\\s\\|\\s(.*)\\s\\|\\s(.*)\\s\\|\\sConnectors.hpp(.*)\\s\\|\\s<json>(.*)")
			// ok:
			//re := regexp.MustCompile("(.*)Connectors\\.hpp(?P<json>(.*))")
			re := extRegexp{regexp.MustCompile(`(.*)Connectors\\.hpp(?P<json>(.*))`)}

			//json := re.FindAllString(s, -1)
			json := re.FindStringSubmatchMap(s)
			if nil == json {
				return "", false
			}
			fmt.Printf("[debug] %s\n", json)
			return json["json"], true
		})
	}

	/*

		re := regexp.MustCompile("(?P<first_char>.)(?P<middle_part>.*)(?P<last_char>.)")
		n1 := re.SubexpNames()
		r2 := re.FindAllStringSubmatch("Super", -1)[0]

		md := map[string]string{}
		for i, n := range r2 {
			fmt.Printf("%d. match='%s'\tname='%s'\n", i, n, n1[i])
			md[n1[i]] = n
		}
		fmt.Printf("The names are  : %v\n", n1)
		fmt.Printf("The matches are: %v\n", r2)
		fmt.Printf("The first character is %s\n", md["first_char"])
		fmt.Printf("The last  character is %s\n", md["last_char"])

		fmt.Println("pcaptool [end]")
	*/
}
