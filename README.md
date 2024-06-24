
# Multi-Backend Cache System

## Overview
The MultiBackend Cache System is designed to facilitate efficient data caching by utilizing multiple caching strategies. 
This system supports Inmemory, Memcache and Redis implementations, allowing for flexible, scalable caching solutions suitable for a variety of applications.
 
## Features
- **Multiple Cache Backends**: Supports Inmemory, Memcache and Redis.
- **Docker Integration**: Easily deployable in a containerized environment with Docker.
- **Concurrency Safe**: Thread-safe implementations ensuring data integrity.
- **Benchmark Tool**: Benchmark tool to test the application performance.
- **Grafana Monitoring Tool**: Interactive and customizable dashboards for real-time data insights.
- **Prometheus**: Monitering and collecting metrics to provide for Grafana.

## Tenent Feature (only for inmemmory)
- **In Memory** - By default, 15% of the system memory will be allotted for the in-memory cache, including the tenant partition. In order to store the cache data into a specific tenant, one has to change config *IsTenantBased* to true and provide *tenantNames* (max 3). The cache memory will be split to each tenant equally. There will be slight changes in Api Endpoints. 
## Table of Contents

1. [Project Structure](#project-structure)
2. [Getting Started](#getting-started)
3. [Configuration](#configuration)
4. [Installation and Deployment](#installation-and-deployment)
5. [Usage](#usage)
6. [Benchmark](#benchmark-tool)

## Project Structure

\`\`\`
multi-backend-cache/
├── docker-compose.yml
├── dockerFile
├── go.mod
├── go.sum
├── prometheus.yml
├── router.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Internal/
│   ├── cache/
│   │   ├── inmemory.go
│   │   ├── interface.go
│   │   ├── memcache.go
│   │   ├── mock_cache.go
│   │   ├── redis.go
│   │   └── tenantCacheManager.go
│   ├── config/
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Handler/
│   │   ├── handler.go
│   │   └── middleware.go
│   └── metrices/
│       └── prometheus_model_info.go
├── kuber-deployment/
│   ├── go-service-deployment.yaml
│   ├── mem-statefulset.yaml
│   ├── redis-configmap.yaml
│   ├── redis-statefulset.yaml
│   ├── twem-config.yaml
│   └── twemproxy.yaml
├── packageUtils/Utils/
│   └── utils.go
└── test/
    ├── inmemory_test.go
    ├── server_test.go
    ├── test_mock.go
    └── test_server_test.go
\`\`\`

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

### Installation and Deployment

1. Clone the repository:
   \`\`\`sh
   git clone https://github.com/FalconDinesh/MultiCacheSystem
   cd multi-backend-cache
   \`\`\`

2. Build the Docker image:
   \`\`\`sh
   docker build -t multi-backend-cache .
   \`\`\`

3. Start the services using Docker Compose:
   \`\`\`sh
   docker-compose up
   \`\`\`

## Configuration

The configuration file \`config.yaml\` is located in the \`Internal/config/\` directory. It includes settings for the different cache backends and other application configurations.


## Usage

The application routes and handlers are defined in \`router.go\` and the \`Internal/Handler/\` directory. API documentation is available in the \`docs/\` directory and can be accessed via [Swagger UI](http://34.234.207.91:8080/swagger/index.html).



## System Options:
- redis
- memcache
- inmemory

## APIs Interact with the cache:
Postman collection is available in the root directory with the following APIs. One can download and import the [collection](https://github.com/sabarivasan007/MultiBackendCacheSystem/blob/main/Multi-Backend-Cache.postman_collection.json) in Postman and test it.

examples:

### Set a cache entry:
#### POST - http://34.234.207.91:8080/cache?system=inmemory
#### with Tenant details  http://34.234.207.91:8080/cache?system=inmemory&tenantID=tenant1
\`\`\`json
{
    "key": "exampleKey", 
    "value": "123",
    "ttl": 100
}
\`\`\`

### Get a cache entry:

GET - http://34.234.207.91:8080/cache/exampleKey?system=inmemory
#### with Tenant details  http://34.234.207.91:8080/cache/tenant1/exampleKey?system=inmemory&tenantID=tenant1

### Delete a cache entry:
DELETE - http://34.234.207.91:8080/cache/dhoni?system=inmemory
#### with Tenant details  http://34.234.207.91:8080/cache?system=inmemory&tenantID=tenant1


### Clear all cache entries:
PUT - http://34.234.207.91:8080/cache/clear?system=inmemory
#### with Tenant details  http://34.234.207.91:8080/cache/clear?system=inmemory&tenantID=tenant1

____

## Test
To test the caches, go to test folder.

* To Test all the cache system at one go
```
go test
```
* To test only the inmemory cache APIs
```
go test inmemory_test.go
```
* To test only the inmemory Tenant cache APIs
```
go test inmemorytenant_test.go
```
* To test only the memcache APIs
```
go test memcache_test.go 
```
* To test only the redis cache APIs
```
 go test redis_test.go 
 ```



## Benchmark Tool
Used to Benchmark the performace of application. It is a separate go application in the directory *[MultiBackendCacheSystem/benchmark-tool/](https://github.com/sabarivasan007/MultiBackendCacheSystem/tree/main/benchmark-tool)*. Detailed overview on how to use the tool is provided in the [README.md](https://github.com/sabarivasan007/MultiBackendCacheSystem/blob/main/benchmark-tool/README.md) file under the *benchmark-tool/* directory.

## Integrating Prometheus and Grafana
 
To monitor the performance of our application, follow these steps to integrate Prometheus and Grafana:

* Grafana 
    * [URL](http://34.234.207.91:3002/d/cdp1812zor0n4a/multicache-monitoring)
    * username - admin
    * password - 28Fo%dT9yb=!wp

### Step 1: Add the Data Source using Prometheus
 
1. Open your Grafana instance in your web browser.
2. Navigate to **Configuration** (the gear icon) > **Data Sources**.
3. Click on **Add data source**.
4. Select **Prometheus** from the list.
5. In the **HTTP** section, enter your Prometheus server URL.
6. Click on **Save & Test** to ensure the connection is working.
 
### Step 2: Import the Dashboard using the Given JSON File

1. In Grafana, go to the **+** icon on the left sidebar and select **Import**.
2. Upload the JSON file provided for the dashboard or paste the JSON content directly. The JSON file is available in the project root directory in the name **grafana-dashboard-config.json**
3. Click **Load**.
4. Select the Prometheus data source that you added in the previous step.
5. Click **Import** to add the dashboard.
 
### Step 3: Activate the Visualization by Running Queries
 
1. Open the imported dashboard.
2. To customize or add new visualizations, click on any panel to edit.
3. Use the **Run Query** button to execute the desired queries and visualize the data.

These steps will enable you to monitor total requests, successful hits, failed hits, overall performance, and throughput of the entire application using Prometheus and Grafana.
