// tqa_pcap_tool project main.go
package main

import (
	"fmt"

	"github.com/frdrolland/pcaptool/cfg"
	"github.com/frdrolland/pcaptool/cli"
	"github.com/frdrolland/pcaptool/parsing"
	"github.com/frdrolland/pcaptool/parsing/ctmetrics"
)

// global variables
var (
	config cfg.Configuration
)

// Verbose prints on standard output string argument if (and only if) -v option has been set in program arguments.
func Verbose(s string, args ...interface{}) {
	if config.Verbose {
		fmt.Printf("%s\n", s)
	}
}

// Main method for pcaptool executable.
func main() {

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	cfg.Global = config
	Verbose("Parsed arguments : %s\n", config)

	// Parse log and extract JSON from each line for specific Regexp
	for _, currFile := range config.Files {
		// element is the element from someSlice for where we are
		Verbose("Importing file %s\n", currFile)
		parsing.ParseLines(currFile, ctmetrics.TryConnectorLine)
	}

}
