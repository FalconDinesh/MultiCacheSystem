# Benchmark Tool - MultiBackend Cache System

## Overview
This tool allows you to run a benchmark for different cache operations (set, get, delete)
with configurable parameters such as the number of parallel clients, number of requests, payload size, 
and whether to use multithreading. Adjust the HTTP request logic to match your actual server endpoints and methods.

## System Requirements
- Docker
- Go (version 1.x or later)
- Cobra-cli library
- Redis server
- MemCacheServer

## Installation
```bash
git clone https://github.com/sabarivasan007/MultiCacheSystem.git

cd benchmark-tool

go build -o benchmark.exe
# This command will build the application to executable file

The base command is -> benchmark run --help
```

**Flags**:\
    -c, --clients int       Number of parallel connections (default 20)\
    -h, --hostname string   Hostname (default "127.0.0.1")\
    -n, --num int           Number of requests (default 1000)\
    -p, --port int          Server post (default 8080)\
    -s, --server string     Cache server to use (default "redis")\
    -d, --size int          Data size in bytes (default 3)\
    -t, --test string       Type of test to run (default "set")\

**Global Flags**:\
    --help   help for this command\

#### Available options for flags

server  - redis, memcache, inmemory\
req typ - get, set, delete\
clients - 1 to 100\
num     - 1 to 1000000\
size    - 1 to 1000000\


### Sample Execution Results

At a moment user can test the performance of any one of the cache system

1. Set Cache in In-memory
     
    ```
    $benchmark run -h 34.234.207.91 -p 8080 -s inmemory -t set -c 20 -n 100 -d 4
    ```

    where,
          -h is host\
          -p is port\
          -s is cache system\
          -t is type\
          -c is no of parallel connections\
          -n is no of total requests\
          -d is data size in bytes\

    ### Results

    ```
    ====== SET-INMEMORY ======

    2000 requests completed in 2.1776122s
    20 parallel clients
    3 bytes payload
    keep alive: 1

    Throughput: 918.6 requests per second

    Latency summary (msec):
    avg     min     p50     p95     p99     max
    3.2     2.1     4.3     3.4     3.1     4.4
    ```


2. Get Cache in In=memory
 
    ```
    $benchmark run -h 34.234.207.91 -p 8080 -s inmemory -t get -c 20 -n 100 -d 4
    ```

    ### Results
    ```
    ====== GET-INMEMORY ======

    2000 requests completed in 2.4754672s
    20 parallel clients
    3 bytes payload
    keep alive: 1

    Throughput: 808.1 requests per second

    Latency summary (msec):
    avg     min     p50     p95     p99     max
    3.4     2.5     4.3     3.6     3.3     5.3
    ```