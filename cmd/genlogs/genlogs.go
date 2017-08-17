package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"

	"github.com/frdrolland/metldr/dto/ctmetrics"
)

var (
	outputFile string
)

const (
	OUTPUT_STRING      = "%s |    OPTIQ[21846-22.10] | notice   | 00 | #9800 |                           Connectors.hpp:185  | %s\n"
	CONNECTOR_TEMPLATE = "{\"sourceType\":\"trading_chain\",\"msgType\":\"system_health_status\",\"data\":{\"optiqSegment\":2,\"optiqSegmentName\":\"ETF\",\"optiqPartitions\":[{\"partitionId\":0,\"partitionNumber\":0,\"serverName\":\"ETF.0\",\"instanceType\":\"primary\",\"publicationTime\":123456789,\"period\":10,\"expectedCoresCount\":12,\"cpuCores\":[{\"core\":32,\"tredzoneTotalLoops\":15742154,\"tredzoneUsedLoops\":1050,\"maxEventsPerLoop\":2103,\"eventsCount\":3943,\"coreUsage_percent\":0.1322562609539235,\"avgEventsPerLoop\":0.000250473982150092},{\"core\":31,\"tredzoneTotalLoops\":17627240,\"tredzoneUsedLoops\":354,\"maxEventsPerLoop\":2103,\"eventsCount\":2640,\"coreUsage_percent\":0.06580566406424109,\"avgEventsPerLoop\":0.00014976819967277917},{\"core\":30,\"tredzoneTotalLoops\":17568205,\"tredzoneUsedLoops\":42,\"maxEventsPerLoop\":2103,\"eventsCount\":2544,\"coreUsage_percent\":0.0730679985278917,\"avgEventsPerLoop\":0.00014480705342406923},{\"core\":29,\"tredzoneTotalLoops\":18325267,\"tredzoneUsedLoops\":8,\"maxEventsPerLoop\":2103,\"eventsCount\":2159,\"coreUsage_percent\":0.06688338955163247,\"avgEventsPerLoop\":0.00011781547302966991},{\"core\":23,\"tredzoneTotalLoops\":5824256,\"tredzoneUsedLoops\":3882877,\"maxEventsPerLoop\":2105,\"eventsCount\":3885581,\"coreUsage_percent\":81.22922444496136,\"avgEventsPerLoop\":0.6671377425717551},{\"core\":28,\"tredzoneTotalLoops\":18388148,\"tredzoneUsedLoops\":8,\"maxEventsPerLoop\":2103,\"eventsCount\":2159,\"coreUsage_percent\":0.06649012682724964,\"avgEventsPerLoop\":0.00011741258554151293},{\"core\":25,\"tredzoneTotalLoops\":9823624,\"tredzoneUsedLoops\":9,\"maxEventsPerLoop\":2103,\"eventsCount\":2861,\"coreUsage_percent\":46.47202776696287,\"avgEventsPerLoop\":0.00029123671671472766},{\"core\":27,\"tredzoneTotalLoops\":18430043,\"tredzoneUsedLoops\":23,\"maxEventsPerLoop\":2103,\"eventsCount\":2504,\"coreUsage_percent\":0.0829728791882216,\"avgEventsPerLoop\":0.00013586511979380624}],\"kafkaUsages\":[{\"topic\":\"OE\",\"partitions\":[]},{\"topic\":\"ME\",\"partitions\":[]}]}]}}"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseCliArgs() (reterr error) {
	usage := `genlogs.

Usage:
  genlogs -n <number> [<output>]
  genlogs -h | --help
  genlogs --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "genlogs tool 1.0", false)

	intervalNano := int64(10)

	if nil == arguments {
		log.Fatal("Missing arguments")
	}

	var number int64
	if nil != arguments["<number>"] {
		snum := arguments["<number>"].(string)
		i, err := strconv.ParseInt(snum, 10, 64)
		if nil != err {
			log.Fatal("Incorrect value for <number> : it must be int64 (%s)", snum)
		}
		number = i
	}

	if nil != arguments["<output>"] {
		outputFile = arguments["<output>"].(string)
	}

	// Reusable buffer
	var buf bytes.Buffer
	buf = bytes.Buffer{}

	var timestamp = time.Now()
	timestamp = timestamp.Add(-1 * time.Duration(number*intervalNano) * time.Second)
	//	fmt.Printf("NOW nano = %d\n", now.UnixNano())

	// Build line protocol message for InfluxDB
	buf.Truncate(0)

	var stat ctmetrics.ConnectorStat
	err := json.Unmarshal([]byte(CONNECTOR_TEMPLATE), &stat)
	if nil != err {
		log.Fatal("ERROR while decoding JSON from file line %s - %s", err)
	}

	stat.SourceType = "trading_chain"
	stat.MsgType = "system_health_status"

	for i := int64(0); i < number; i++ {
		if i > 0 {
			timestamp = timestamp.Add(time.Duration(intervalNano) * time.Second)
		}
		//
		for idxPart, _ := range stat.Data.OptiqPartitions {

			partStat := &stat.Data.OptiqPartitions[idxPart]

			partStat.PublicationTime = timestamp.UnixNano()

			for _, coreStat := range partStat.CPUCores {
				// Reinit buffer
				coreStat.EventsCount = 0
			}
		}
		b, err := json.Marshal(stat)
		if nil != err {
			log.Fatal("Could not encode in JSON : %s", err)
		}
		fmt.Printf(OUTPUT_STRING, timestamp.Format("2006/01/02 15:04:05.999999"), b)
	}

	return
}

// Tool to generate large log files from a template, to test volumes/performances.
func main() {
	ParseCliArgs()
}
