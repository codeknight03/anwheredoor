package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codeknight03/anywheredoor/pkg/config"
	"github.com/codeknight03/anywheredoor/pkg/server"
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

	proxy := server.NewReverseProxy(rpcfg)

	log.Printf("Starting reverse proxy on port %s...", rpcfg.ListenPort)
	if err := http.ListenAndServe(":"+rpcfg.ListenPort, proxy); err != nil {
		log.Printf("Last config parsed: %v", rpcfg)
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Print(rpcfg)

}
