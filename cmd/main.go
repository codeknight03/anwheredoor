package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/codeknight03/anywheredoor/pkg/config"
	"github.com/codeknight03/anywheredoor/pkg/server"
)

//Question: 	 How do we add config Reload?
//Rationale:	 Binding the running server to a unix socket
//				 and listen to that socket for update config.
//Alternatives:  1. Watch Files ( Overhead was high, also unnecessary
// 				 reloads may happen as reload is automatic)
//				 2. Reload Endpoint  (The end point will have to be
// 				 secured as anyone can re-route traffic and again
// 				 that is too much overhead)
//TODO: Implement a HTTP based reload endpoint with security.

const (
	SOCKET_PATH = "/tmp/anywheredoor.sock"
	LOG_LEVEL   = "DEBUG"
)

type ConfigUpdateMessage struct {
	Config []byte `json:"config"`
}

func setupLogger(envLogLevel string) {
	var level slog.Level

	switch {
	case envLogLevel == slog.LevelDebug.String():
		level = slog.LevelDebug
	case envLogLevel == slog.LevelInfo.String():
		level = slog.LevelInfo
	case envLogLevel == slog.LevelWarn.String():
		level = slog.LevelWarn
	case envLogLevel == slog.LevelError.String():
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

}

func main() {
	//logging setup
	setupLogger(LOG_LEVEL)

	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	start := flag.Bool("start", false, "Start the reverse proxy")
	update := flag.Bool("update", false, "Update the configuration")
	flag.Parse()

	if !*start && !*update {
		fmt.Println("Usage: anywheredoor -start|-update -c config.yaml")
		os.Exit(1)
	}

	if *start {
		// Remove existing socket file if it exists
		os.Remove(SOCKET_PATH)

		byteData, err := os.ReadFile(*configPath)
		if err != nil {
			slog.Error("Cannot read config file", "error", err)
			os.Exit(1)
		}

		rpcfg, err := config.RpConfigFromBytes(byteData)
		if err != nil {
			slog.Error("Cannot unmarshal the config", "error", err)
			os.Exit(1)
		}

		proxy := server.NewReverseProxy(rpcfg)

		listener, err := net.Listen("unix", SOCKET_PATH)
		if err != nil {
			slog.Error("Failed to create socket", "error", err)
			os.Exit(1)
		}
		defer listener.Close()

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					slog.Warn("Error accepting connection to update config.", "error", err)
					continue
				}

				go handleConfigUpdate(conn, proxy)
			}
		}()

		slog.Info("Starting reverse proxy.", "port", rpcfg.ListenPort)
		if err := http.ListenAndServe(":"+rpcfg.ListenPort, proxy); err != nil {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}

	} else if *update {

		conn, err := net.Dial("unix", SOCKET_PATH)
		if err != nil {
			slog.Error("Failed to connect to running proxy.", "error", err)
			os.Exit(1)
		}
		defer conn.Close()

		// Read and send new configuration
		byteData, err := os.ReadFile(*configPath)
		if err != nil {
			log.Fatalf("Cannot read config file: %s", err)
		}

		msg := ConfigUpdateMessage{
			Config: byteData,
		}
		//Inter Process Communication using json.
		//TODO: Implement this with protobuf for learning .proto.
		if err := json.NewEncoder(conn).Encode(msg); err != nil {
			slog.Error("Failed to send config update.", "error", err)
			os.Exit(1)
		}

		slog.Debug("Configuration update sent successfully.")
	}
}

func handleConfigUpdate(conn net.Conn, proxy *server.ReverseProxy) {
	defer conn.Close()

	var msg ConfigUpdateMessage
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		slog.Warn("Error decoding config update.", "error", err)
		return
	}

	rpcfg, err := config.RpConfigFromBytes(msg.Config)
	if err != nil {
		slog.Warn("Error parsing new config.", "error", err)
		return
	}

	proxy.UpdateConfig(rpcfg)
	slog.Info("Configuration updated successfully")
}
