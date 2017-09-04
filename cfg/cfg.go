package cfg

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	Verbose bool
	Files   []string
	Command string
	Source  string
	Input   struct {
		Kafka struct {
			Brokers string `json:"brokers"`
			Group   string `json:"group"`
			Topic   string `json:"topic"`
		}
	}
	Output struct {
		InfluxDB struct {
			Protocol  string
			Url       string
			UdpAddr   string
			Database  string
			RetPolicy string
			User      string
			Password  string
		}
	}
}

var (
	Global Configuration
)

func Init() {
	viper.SetConfigName("metldr") // name of config file (without extension)
	//viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")

	viper.AddConfigPath("/etc/metldr")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.metldr") // call multiple times to add many search paths
	viper.AddConfigPath("./config")      // optionally look for config in the working directory

	viper.SetDefault("logging.dir", "./logs")
	viper.SetDefault("input.kafka.brokers", "127.0.0.1:9092")
	viper.SetDefault("output.influxdb.url", "http://localhost:8086")
	viper.SetDefault("output.influxdb.database", "ct")
	viper.SetDefault("output.influxdb.udpaddr", "127.0.0.1:8089")
	viper.SetDefault("output.influxdb.retpolicy", "autogen")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&Global)
	if err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//TODO utiliser un log
		fmt.Println("Config file changed:", e.Name)
	})

	fmt.Printf("%s = %s\n", "input.kafka.brokers", Global.Input.Kafka.Brokers)
	fmt.Printf("%s = %s\n", "output.influxdb.udpaddr", Global.Output.InfluxDB.UdpAddr)
	fmt.Printf("%s = %s\n", "output.influxdb.database", Global.Output.InfluxDB.Database)
	fmt.Printf("%s = %s\n", "output.influxdb.retpolicy", Global.Output.InfluxDB.RetPolicy)
	fmt.Printf("%s = %s\n", "output.influxdb.user", Global.Output.InfluxDB.User)
	fmt.Printf("%s = %s\n", "output.influxdb.password", Global.Output.InfluxDB.Password)
	fmt.Println(viper.AllSettings())
}
