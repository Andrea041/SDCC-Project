# Implementation of 2 distributed leader election algorithms for Distributed System & Cloud Computing course

## Introduction
A distributed system is a collection of independent network nodes that do not share memory. Each processor on every node has its own memory and communicates with others through the network. Network communication is implemented via a process on one machine communicating with a process on another machine.

Many algorithms used in distributed systems require the presence of a coordinator or leader responsible for carrying out necessary functions for other processes in the system. For the selection of a new coordinator, election algorithms have been specifically designed.

Election algorithms select a process from a group of processes to act as coordinator. If the leader stops for any reason, a new one is elected. Election obviously requires achieving distributed consensus among the nodes in the system.

The algorithms are:
- Bully Algorithm
- Chang-Roberts Algorithm

## Requirements & Instruction to run the code
### Software requirements
Software that must be installed on your computer:
- go 1.22.2 darwin/amd64
- Docker 25.0.3, build 4debf41
- Docker Compose v2.24.6-desktop.1

### Installation
To install this application on your computer, run this git command:
```bash
git clone https://github.com/Andrea041/Distributed-System-And-Cloud-Computing-Project
```

### Run code
