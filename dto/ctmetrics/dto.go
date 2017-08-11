// Parsing of ct-metrics lgos (metrics generated by Optiq core-trading components).
package ctmetrics

// Bind stats of "Connectors" log lines.
// Connectors.hpp:185
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

type MonitorWorkerAvail struct {
	// MonitorWorkerAvail.hpp:321
	TCPSourceAvail struct {
		Time  string `json:"time"`
		Name  string `json:"name"`
		State string `json:"state"`
		Core  int    `json:"core"`
		ID    int    `json:"id"`
		IP    string `json:"ip"`
		Port  string `json:"port"`
		ErrNo int    `json:"errNo"`
	} `json:"tcpSourceAvail"`
	// MonitorWorkerAvail.hpp:354
	MulticastEmitterAvail struct {
		Time      string `json:"time"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Core      int    `json:"core"`
		ID        int    `json:"id"`
		Interface string `json:"interface"`
		ErrNo     int    `json:"errNo"`
	} `json:"multicastEmitterAvail"`
	// MonitorWorkerAvail.hpp:436
	PePersistenceAvail struct {
		Time        string `json:"time"`
		Name        string `json:"name"`
		State       string `json:"state"`
		Core        int    `json:"core"`
		ID          int    `json:"id"`
		Brokers     string `json:"brokers"`
		NbBrokers   int    `json:"nbBrokers"`
		NbBrokersUp int    `json:"nbBrokersUp"`
		Queues      int    `json:"queues"`
		ErrCode     int    `json:"errCode"`
	} `json:"pePersistenceAvail"`
	// MonitorWorkerAvail.hpp:525
	MdgAvail struct {
		Time         string `json:"time"`
		MdgID        int    `json:"mdgId"`
		OptiqSegment int    `json:"optiqSegment"`
		PartitionID  int    `json:"partitionId"`
		MdgState     string `json:"mdgState"`
	} `json:"mdgAvail"`
}

// MonitorWorkerCount.hpp:221
type MonitorWorkerCount struct {
	TCPSourceCount struct {
		Time                  string  `json:"time"`
		Name                  string  `json:"name"`
		Core                  int     `json:"core"`
		ID                    int     `json:"id"`
		CumulatedConnections  int     `json:"cumulatedConnections"`
		CumulatedUserMessages int     `json:"cumulatedUserMessages"`
		CumulatedBadMessages  int     `json:"cumulatedBadMessages"`
		RateElapsed           float64 `json:"rateElapsed"`
		RateUserMessages      float64 `json:"rateUserMessages"`
		RateBadMessages       float64 `json:"rateBadMessages"`
	} `json:"tcpSourceCount"`
	AggregCount struct {
		Time                  string  `json:"time"`
		Name                  string  `json:"name"`
		Core                  int     `json:"core"`
		ID                    int     `json:"id"`
		CumulatedUserMessages int     `json:"cumulatedUserMessages"`
		CumulatedBadMessages  int     `json:"cumulatedBadMessages"`
		RateElapsed           float64 `json:"rateElapsed"`
		RateUserMessages      float64 `json:"rateUserMessages"`
		RateBadMessages       float64 `json:"rateBadMessages"`
	} `json:"aggregCount"`
}