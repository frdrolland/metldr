package cli

import (
	"github.com/docopt/docopt-go"

	"github.com/frdrolland/metldr/cfg"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseCliArgs() (config cfg.Configuration, err error) {
	usage := `metldr.

Usage:
  metldr show file <filename>...
  metldr show  kafka <group> <topic> [--brokers <brokerlist>]
  metldr import file <filename>...
  metldr import  kafka <group> <topic> [--brokers <brokerlist>]
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

	if arguments["kafka"].(bool) {
		config.Source = "kafka"
		config.Input.Kafka.Group = arguments["<group>"].(string)
		config.Input.Kafka.Topic = arguments["<topic>"].(string)
		brokers := arguments["<brokerlist>"]
		if nil != brokers && "" != brokers.(string) {
			config.Input.Kafka.Brokers = brokers.(string)
		}
	} else if arguments["file"].(bool) {
		config.Source = "file"
	}

	// Finally set global config to read item
	return config, nil
}
