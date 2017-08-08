package cli

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/frdrolland/pcaptool/cfg"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseCliArgs() (config cfg.Configuration, err error) {
	usage := `pcaptool.

Usage:
  pcaptool import <filename>...
  pcaptool -h | --help
  pcaptool --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	config = cfg.Configuration{}

	arguments, _ := docopt.Parse(usage, nil, true, "Pcap Tool 1.0", false)
	config.Files = arguments["<filename>"].([]string)
	fmt.Printf("Parsed arguments : %s\n", config)

	return
}
