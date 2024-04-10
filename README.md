# Implementation of two distributed leader election algorithms for Distributed System & Cloud Computing course

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
- go 1.22.2 darwin/amd64 (or any compatible version)
- Docker 25.0.3 (build 4debf41 recommended)
- Docker Compose v2.24.6-desktop.1 or later

### Choose algorithm to test
To choose the algorithm you have to set "true" or "false" in config.json file:
```json
"algorithm": {
    "ChangAndRobert": "false",
    "Bully": "true"
  }
```


### Run code without AWS EC2 istance
If you want try this application without AWS EC2 istance, you have to clone this repository on your computer:
```bash
git clone https://github.com/Andrea041/Distributed-System-And-Cloud-Computing-Project
```
Then start Docker daemon and run Docker Compose:
```bash
sudo service docker start
docker-compose up -d
```

### Run code with AWS EC2 istance
If you want try this application with AWS EC2 istance (Amazon Linux OS), you have to create an EC2 istance and connect to it:
```bash
ssh -i <path_to_PEM> ec2-user@<ip-EC2-instance>
```
First of all you have to install git command:
```bash
sudo yum update -y
sudo yum install git -y
```
Then you have to clone this repository:
```bash
git clone https://github.com/Andrea041/Distributed-System-And-Cloud-Computing-Project
```
Now install Docker:
```bash
sudo yum install -y docker
```
And Docker Compose:
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```
Run Docker daemon:
```bash
sudo service docker start
```
Finally run Docker Compose command:
```bash
sudo docker-compose -f compose.yaml up
```
