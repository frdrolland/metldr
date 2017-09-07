package input

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"

	"github.com/frdrolland/metldr/cfg"
	"github.com/frdrolland/metldr/ctmetrics"
)

var (
	partitions = "all"
	offset     = "newest"
	logger     = log.New(os.Stderr, "", log.LstdFlags)
)

//TODO A Déplacer !!!
func printErrorAndExit(code int, format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

//
//
//
func ConsumeKafkaMetrics() {

	broker := cfg.Global.Input.Kafka.Brokers
	group := cfg.Global.Input.Kafka.Group
	topics := []string{cfg.Global.Input.Kafka.Topic}

	// Init config
	config := cluster.NewConfig()
	sarama.Logger = logger
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	switch offset {
	case "oldest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		log.Fatal("-offset should be `oldest` or `newest`")
	}

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(strings.Split(broker, ","), group, topics, config)
	if err != nil {
		printErrorAndExit(69, "Failed to start consumer: %s", err)
	}
	defer consumer.Close()

	// Create signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// Consume all channels, wait for signal to exit
	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
				ctStat := ctmetrics.GetStat(string(msg.Value))
				ctmetrics.ProcessEvent(ctStat)
				consumer.MarkOffset(msg, "")
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				logger.Printf("Rebalanced: %+v\n", ntf)
			}
		case err, more := <-consumer.Errors():
			if more {
				logger.Printf("Error: %s\n", err.Error())
			}
		case <-sigchan:
			return
		}
	}
}
