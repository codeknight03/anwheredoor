#!/bin/bash

# Start demo services
start_containers() {
  echo "Starting demo HTTP servers..."
  docker run --rm -d -p 9001:80 --name server1 -e SERVER_NAME="Server 1" hashicorp/http-echo -listen=:80 -text="Hello from Server 1"
  docker run --rm -d -p 9002:80 --name server2 -e SERVER_NAME="Server 2" hashicorp/http-echo -listen=:80 -text="Hello from Server 2"
  docker run --rm -d -p 9003:80 --name server3 -e SERVER_NAME="Server 3" hashicorp/http-echo -listen=:80 -text="Hello from Server 3"
  echo "Servers are running at:"
  echo " - http://localhost:9001"
  echo " - http://localhost:9002"
  echo " - http://localhost:9003"
}

# Stop demo services
stop_containers() {
  echo "Stopping demo HTTP servers..."
  docker stop server1 server2 server3
}

# Display usage
usage() {
  echo "Usage: $0 {start|stop}"
  exit 1
}

# Handle user input
case "$1" in
  start) start_containers ;;
  stop) stop_containers ;;
  *) usage ;;
esac

