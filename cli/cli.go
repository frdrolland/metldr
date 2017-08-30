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
  metldr show  kafka <brokerlist> <group> <topic>
  metldr import file <filename>...
  metldr import kafka <brokerlist> <group> <topic>
  metldr -h | --help
  metldr --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	config = cfg.Configuration{}

	arguments, _ := docopt.Parse(usage, nil, true, "Pcap Tool 1.0", false)
	config.Files = arguments["<filename>"].([]string)
	//config.BrokerList = arguments["<brokerlist>"].([]string)
	//fmt.Printf("ARGS = %s\n", arguments)

	if arguments["import"].(bool) {
		config.Command = "import"
	} else if arguments["show"].(bool) {
		config.Command = "show"
	}

	if arguments["kafka"].(bool) {
		config.Source = "kafka"
		config.BrokerList = arguments["<brokerlist>"].(string)
		config.KafkaGroup = arguments["<group>"].(string)
		config.KafkaTopic = arguments["<topic>"].(string)
	} else if arguments["file"].(bool) {
		config.Source = "file"
	}

	// Finally set global congig to read item
	return config, nil
}
