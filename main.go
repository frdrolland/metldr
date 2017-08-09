// tqa_pcap_tool project main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Command-line arguments parsing
	config, _ := cli.ParseCliArgs()
	Verbose("Parsed arguments : %s\n", config)

	// Parse log and extract JSON from each line for specific Regexp
	for _, currFile := range config.Files {
		// element is the element from someSlice for where we are
		Verbose("Importing file %s\n", currFile)
		parsing.ParseLines(currFile, func(s string) (string, bool) {

			var buf bytes.Buffer
			buf = bytes.Buffer{}

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

				// Build line protocol message for InfluxDB
				buf.Truncate(0)

				//TODO Code à optimiser: (remplacer les fmt.Sprint par des buf.Write 'simples')
				for _, partStat := range newStat.Data.OptiqPartitions {

					for _, coreStat := range partStat.CPUCores {
						// Reinit buffer

						// measurement
						buf.WriteString("system_health_status")

						// tagset
						buf.WriteString(",")
						buf.WriteString(fmt.Sprintf(`part_id=%d,part_num=%d,server_name=%s,type=%s,core=%d`, partStat.PartitionID, partStat.PartitionNumber, partStat.ServerName, partStat.InstanceType, coreStat.Core))

						// tagset
						buf.WriteString(" ")
						buf.WriteString(fmt.Sprintf(`tz_loops_total=%d,tz_loops_used=%d,events=%d,core_usage_pct="%f",avg_events_per_loop="%f",max_events_per_loop=%d`, coreStat.TredzoneTotalLoops, coreStat.TredzoneUsedLoops, coreStat.EventsCount, coreStat.CoreUsagePercent, coreStat.AvgEventsPerLoop, coreStat.MaxEventsPerLoop))

						// timestamp
						buf.WriteString(" ")
						buf.WriteString(fmt.Sprintf("%d", partStat.PublicationTime))

						buf.WriteString("\n")

					}
				}

				switch command := config.Command; command {
				case "import":
					// Import data ni InfluxDB
					resp, err := http.Post("http://localhost:8086/write?db=testfro", "text/plain", &buf)
					if nil != err {
						fmt.Printf("ERROR while uploading on InfluxDB: %s\n", err)
					} else {
						fmt.Printf("UPLOADED: %s - STATUS=%d\n", buf.String(), resp.Status)
					}
				case "show":
					// Show only generated data on standard output
					fmt.Printf("%s", buf.String())
				default:
					log.Fatal(fmt.Sprintf("Unknown command: %s", command))
					os.Exit(10)
				}

			}

			return extracted, true
		})
	}

}
