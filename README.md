# Distributed URL Shortener

## Overview

This project implements a Distributed URL Shortener using the Chord protocol.
URL shortener service that can store and retrieve URLs in a distributed manner.
Each node in the DHT is represented by a port running on the system.
Each node has a Finger table, maintaining node IDs
On receiving a URL shortening request, the URL shortener service creates the hash key of the log URL and returns the response to the client
The generated hash key in the above step has a key ID and chord protocol determines the responsible node to store the key based on the key ID.
During URL retrieval from the short key,  using chord protocol, the request will be hopped between eligible nodes to find the responsible node 
During each node joins and leaving, finger table entries are stabilized and replication happens to replicate data to the successor and Predecessor.


## Table of Contents

- [Team Members and Division of Work](#team-members-and-division-of-work)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Setup](#setup)
- [Usage](#usage)
  - [Store a URL](#store-a-url)
  - [Lookup a URL](#lookup-a-url)
  - [Leave the Network](#leave-the-network)
- [Troubleshooting](#troubleshooting)
- [Conclusion](#conclusion)

## Team Members and Division of Work

- **Nandini**: Implemented the core DHT functionality including node joining, finger table management, and stabilization, and prepared a readme for the project.
- **Rahul**: Developed the URL shortener service, integrated it with the DHT, and prepared a presentation for demo.
- **Yashashvi**: Handled HTTP API setup, API endpoint creation, and implemented the leave functionality for nodes and prepared the report for the project

## Prerequisites

- Go 1.21 or higher installed on your system from https://go.dev/doc/install based on your OS
- After installing, please make sure to export the path based on your OS and source file
  Eg:
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
- Git installed on your system from https://git-scm.com/book/en/v2/Getting-Started-Installing-Git.

## Installation

### Step 1: Clone the Repository

First, clone the project repository from GitHub.

```sh
git clone [<repository-url>](https://github.com/RahulDhiman93/DS_Final_Project/
cd DS_Final_Project
```

### Step 2: Initialize and Start Nodes

This project starts with multiple servers using multiple ports, each representing a node in the DHT.
By default, the project is set up to start 5 nodes.

### Step 3: Compile and Run the Project

To compile and run the project, follow these steps:

1. Navigate to the project directory.
2. Run the following command to compile and start the servers:

```sh
go run main.go
```

This command will start 5 HTTP servers on ports 8000, 8001, 8002, 8003, and 8004.

## Setup

### Configuring Nodes

The project is designed to start nodes on predefined ports. If you need to start nodes on different ports or add more nodes, modify the `nodes` slice in the `main` function of the `main.go` file.

```go
nodes := []int{8000, 8001, 8002, 8003, 8004}
```

### Starting Nodes

To start the nodes, run the `main.go` file:

```sh
go run main.go
```

This command will start all the nodes defined in the `nodes` slice. Each node will automatically join the DHT network and initialize its finger table.

## Usage

### Store a URL

To store a URL, send a POST request to any node's `/store` endpoint.

```sh
curl -X POST http://127.0.0.1:8000/store -H "Content-Type: application/json" -d '{"url":"http://example.com"}'
```

### Lookup a URL

To lookup a URL, send a POST request to any node's `/lookup` endpoint with the short key.

```sh
curl -X POST http://127.0.0.1:8000/lookup -H "Content-Type: application/json" -d '{"key":"<shortKey>"}'
```

### Leave the Network

To make a node leave the network, send a POST request to the node's `/leave` endpoint with the port number.

```sh
curl -X POST http://127.0.0.1:8000/leave -H "Content-Type: application/json" -d '{"port":8000}'
```

## Troubleshooting

- **Port already in use**: If a server fails to start, check if the port is already in use. Ensure no other services are running on the specified ports (8000-8004) before starting the servers.
- **Network stabilization**: After a node leaves, the remaining nodes will stabilize and update their finger tables to reflect the changes in the network. If nodes do not stabilize properly, ensure all nodes are started in the correct sequence.

## Conclusion

This project demonstrates a basic implementation of a DHT using the Chord protocol and provides a practical application in the form of a URL shortener service.

---

This Readme provides a comprehensive overview of the project, detailed installation and setup instructions, usage examples, and troubleshooting tips.
If you encounter any issues or have questions, please reach out to the project contributors for assistance.
