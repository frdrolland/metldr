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
	viper.SetDefault("output.influxdb.url", "http://localhost:8086")
	viper.SetDefault("output.influxdb.database", "ct")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//TODO utiliser un log
		fmt.Println("Config file changed:", e.Name)
	})

	dburl := viper.Get("output.influxdb.url")
	dbname := viper.Get("output.influxdb.database")

	fmt.Printf("%s = %s\n", "output.influxdb.url", dburl)
	fmt.Printf("%s = %s\n", "output.influxdb.database", dbname)
	fmt.Printf("%s = %s\n", "output.influxdb.user", dbname)
	fmt.Printf("%s = %s\n", "output.influxdb.password", dbname)
	fmt.Println(viper.AllSettings())

}
