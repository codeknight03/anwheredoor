package main

import (
	"fmt"
	"os"

	"github.com/codeknight03/anywheredoor/pkg/config"
)

func main() {

	byteData, err := os.ReadFile("/home/codeknight/anywheredoor/config.yaml")
	if err != nil {
		fmt.Printf("Cannot read config file: %s", err)
	}

	rpcfg, err := config.RpConfigFromBytes(byteData)
	if err != nil {
		fmt.Printf("Cannot unmarshal the config: %s", err)
	}

	fmt.Print(rpcfg)

}
