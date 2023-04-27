## Request counter: A Persistent Go HTTP Server
This repository contains a Go HTTP server that responds with a counter of the total number of requests it has received during the previous 60 seconds (moving window). The server persists the request count data to a file, ensuring that the correct numbers are returned even after restarting the server.

## Features
1. Go HTTP server using the standard library
2. 60-second moving window request counter
3. Data persistence to a file
4. Graceful shutdown
5. Race condition tests

## Getting Started
### Prerequisites
Go 1.17 or higher

#### Installation and Usage
Clone the repository:
```bash
git clone https://github.com/mhsnrafi/request-counter.git
```
Change the directory:
```bash
cd request-counter
```

- To build the project, run: 
```
make build
```
- To run the server:
```cd output
./server
```

- To run the race test, execute the following command:
```
make test
```
- To clean the build artifacts, run: 
```
make clean
```


**The server will start listening on localhost:8080.**

Send a request to the server:
```http request
curl http://localhost:8080
```
The server will respond with the number of requests received during the previous 60 seconds (moving window).

## Code Explaination
- counter.go: This file defines the Counter structure and its associated methods. The Counter structure is a thread-safe counter that only counts events within a specified time window. It also handles data persistence to a file.

- server.go: This file defines the Run() function that starts the HTTP server with the specified Counter instance. It sets up the request handler, which increments the counter and responds with the total number of requests in the last 60 seconds. This file also contains functions to start the server, handle OS signals (like SIGINT and SIGTERM), and gracefully shut down the server.

- counter_test.go: counter_test.go: This file contains tests for the Counter, including tests for the Increment, Count, Save, Load, FindFirstIndexInWindow, and CleanOldTimestamps functionalities.

- server_test.go: This file contains tests for the server, including a test to check the Increment and Count functionalities of the Counter.
