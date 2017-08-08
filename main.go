// tqa_pcap_tool project main.go
package main

import (
	"fmt"

	"github.com/frdrolland/pcaptool/cli"
)

func main() {

	config, _ := cli.ParseCliArgs()

	fmt.Printf("config = %s\n", config)

	fmt.Println("pcaptool [end]")

}
