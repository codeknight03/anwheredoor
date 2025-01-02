# Basic Reverse Proxy Example 

1. Pre-Reuisites:
- A Unix like shell to execute the script.
- Docker Daemon Running

2. Run the Script to run 3 servers:

```
./three-servers.sh 
```

3. Open the configured urls in your browser: 

``` 
http://localhost:8080/ 
# Returns Hello from Server 1

http://localhost:8080/second
# Returns Hello from Server 2

http://localhost:8080/third
# Returns Hello from Server 3
```
