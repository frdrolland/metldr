package process

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frdrolland/metldr/cfg"
	"github.com/frdrolland/metldr/dto/ctmetrics"
)

//
// Processes an event and send it to stdout or to InfluxDB, depending on which command is executed.
//
func ProcessEvent(newStat ctmetrics.ConnectorStat) error {

	var buf bytes.Buffer
	buf = bytes.Buffer{}

	// Build line protocol message for InfluxDB
	buf.Truncate(0)

	//TODO Code à optimiser: (remplacer les fmt.Sprint par des buf.Write 'simples')
	for _, partStat := range newStat.Data.OptiqPartitions {

		for _, coreStat := range partStat.CPUCores {
			// Reinit buffer

			// measurement
			buf.WriteString("connector")

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

	switch command := cfg.Global.Command; command {
	case "import":
		// Import data ni InfluxDB
		resp, err := http.Post("http://localhost:8086/write?db=testfro", "text/plain", &buf)
		if nil != err {
			fmt.Printf("ERROR while uploading on InfluxDB: %s\n", err)
			return err
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
	return nil

}
