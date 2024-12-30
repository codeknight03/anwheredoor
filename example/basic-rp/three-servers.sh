#!/bin/bash

# Start demo services
start_containers() {
  echo "Starting demo HTTP servers..."
  docker run --rm -d -p 9001:80 --name server1 kennethreitz/httpbin
  docker run --rm -d -p 9002:80 --name server2 kennethreitz/httpbin
  docker run --rm -d -p 9003:80 --name server3 kennethreitz/httpbin
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

