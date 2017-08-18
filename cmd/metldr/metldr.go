// tqa_pcap_tool project main.go
package main

import (
	"fmt"

	"github.com/frdrolland/metldr/cfg"
	"github.com/frdrolland/metldr/cli"
	"github.com/frdrolland/metldr/ctmetrics"
	"github.com/frdrolland/metldr/parsing"
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

// Main method for metldr executable.
func main() {

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	cfg.Global = config
	Verbose("Parsed arguments : %s\n", config)

	cfg.Init()

	// Parse log and extract JSON from each line for specific Regexp
	for _, currFile := range config.Files {
		// element is the element from someSlice for where we are
		Verbose("Importing file %s\n", currFile)
		parsing.ParseLines(currFile, ctmetrics.ParseConnectorLines)
	}

}
