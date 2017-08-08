// tqa_pcap_tool project main.go
package main

import (
	"encoding/json"
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

type ConnectorStat struct {
	Data struct {
		OptiqPartitions []struct {
			CPUCores []struct {
				AvgEventsPerLoop   float64 `json:"avgEventsPerLoop"`
				Core               int     `json:"core"`
				CoreUsagePercent   float64 `json:"coreUsage_percent"`
				EventsCount        int     `json:"eventsCount"`
				MaxEventsPerLoop   int     `json:"maxEventsPerLoop"`
				TredzoneTotalLoops int     `json:"tredzoneTotalLoops"`
				TredzoneUsedLoops  int     `json:"tredzoneUsedLoops"`
			} `json:"cpuCores"`
			ExpectedCoresCount int    `json:"expectedCoresCount"`
			InstanceType       string `json:"instanceType"`
			KafkaUsages        []struct {
				Partitions []interface{} `json:"partitions"`
				Topic      string        `json:"topic"`
			} `json:"kafkaUsages"`
			PartitionID     int    `json:"partitionId"`
			PartitionNumber int    `json:"partitionNumber"`
			Period          int    `json:"period"`
			PublicationTime int    `json:"publicationTime"`
			ServerName      string `json:"serverName"`
		} `json:"optiqPartitions"`
		OptiqSegment     int    `json:"optiqSegment"`
		OptiqSegmentName string `json:"optiqSegmentName"`
	} `json:"data"`
	MsgType    string `json:"msgType"`
	SourceType string `json:"sourceType"`
}

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

func main() {
	//t := time.Unix(0, 1501681526043505000)

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	Verbose("Parsed arguments : %s\n", config)

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

			extracted := md["json"]

			if "" == extracted {
				return "", false
			}

			newStat := ConnectorStat{}
			if extracted != "" {
				err := json.Unmarshal([]byte(extracted), &newStat)
				if nil != err {
					fmt.Printf("ERROR while decoding JSON from file line %s - %s", extracted, err)
				}
				fmt.Printf("JSON STRUC ==== %s\n", newStat.Data.OptiqPartitions[0].PublicationTime)
			}

			return extracted, true
		})
	}

}
