package cli

import (
	"github.com/docopt/docopt-go"

	"github.com/frdrolland/metldr/cfg"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseCliArgs() (config cfg.Configuration, err error) {
	usage := `metldr.

Usage:
  metldr show <filename>...
  metldr import <filename>...
  metldr -h | --help
  metldr --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	config = cfg.Configuration{}

	arguments, _ := docopt.Parse(usage, nil, true, "Pcap Tool 1.0", false)
	config.Files = arguments["<filename>"].([]string)

	if arguments["import"].(bool) {
		config.Command = "import"
	} else if arguments["show"].(bool) {
		config.Command = "show"
	}

	return
}
