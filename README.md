# Go Web Framework Performance Benchmark

English | [中文](README_CN.md)

This project provides a simple performance test suite to compare popular Go web frameworks:
- Hertz (High-performance HTTP framework by CloudWeGo)
- Gin (Popular lightweight HTTP framework)
- Go-Zero (Microservice framework)
- Kratos (Microservice framework by Bilibili)
- net/http (Go standard library)

This project was generated using Cursor
Testing tools used:
- [hey](https://github.com/rakyll/hey)
- [wrk](https://github.com/wg/wrk) as benchmark tools

## Usage
- Deploy to docker and run tests
```docker compose up -d```
- Export results to local
```docker cp go-http-benchmark-benchmark-1:/app/benchmark_results/summary_20250419_022006.md .```

## Cursor Prompt
- Cursor Rules settings are attached at the end.
- During project generation, input the following prompt:
> **Create a performance test plan for Hertz, Gin, Go-zero, Kratos and net/http frameworks**
- Output the following content:
> I'll create a complete Go project to test the performance of Hertz, Gin, Go-zero, Kratos and standard net/http. This project will fairly compare the performance of these frameworks in different scenarios. First create project structure:
> ... (remaining content follows same pattern of translation)

## Test Scenarios

This benchmark covers the following scenarios:

1. **Simple Response** (`/hello`): Returns simple Hello World response
2. **JSON Response** (`/json`): Returns JSON formatted data list
3. **Path Parameter** (`/params/:id`): Processes and returns URL path parameter
4. **Query Parameter** (`/query?name=xxx`): Processes and returns URL query parameter
5. **POST Request** (`/post`): Processes JSON formatted POST request body

## Project Structure

```
.
├── cmd
│   ├── benchmark/       # Benchmark execution tool
│   └── server/          # Framework server implementations
│       ├── gin/
│       ├── gozero/
│       ├── hertz/
│       ├── kratos/
│       └── nethttp/
├── docker/              # Docker config files
├── internal/
│   └── benchmark/       # Benchmark tools and shared code
├── docker-compose.yml   # Docker Compose config
├── go.mod               # Go module definition
└── README.md            # Project description
└── README_CN.md            # Project description
```

## Running Tests

### Prerequisites

- Go 1.20 or higher
- Docker and Docker Compose (optional, for containerized testing)
- hey or wrk (HTTP benchmarking tools)

### Install Testing Tools

```bash
# Install hey tool
go install github.com/rakyll/hey@latest

# Or install wrk (Linux/Mac)
# Ubuntu/Debian
apt-get install wrk
# Mac
brew install wrk
```

### Local Testing
1. Start servers for each framework (run each in a separate terminal):
```bash
# Start standard HTTP service
go run cmd/server/nethttp/main.go

# Start Gin service
go run cmd/server/gin/main.go

# Start Hertz service
go run cmd/server/hertz/main.go

# Start Go-Zero service
go run cmd/server/gozero/main.go

# Start Kratos service
go run cmd/server/kratos/main.go
```

2. Run benchmark

```bash
# Test all frameworks
go run cmd/benchmark/main.go

# Test specific framework
go run cmd/benchmark/main.go -framework gin

# Custom test parameters
go run cmd/benchmark/main.go -framework hertz -n 10000 -c 200 -d 30s -tool wrk
```

### Using Docker Compose

```
# Build and start all services
docker-compose up -d

# Run benchmarks
docker-compose run benchmark
```

## Command Line Parameters
The benchmark tool supports these parameters:

- `framework`: Framework to test (nethttp, gin, hertz, gozero, kratos, all)
- `n`: Total requests
- `c`: Concurrent connections
- `d`: Test duration (e.g. "10s", "1m")
- `tool`: Benchmark tool (hey, wrk)
- `format`: Report format (json, csv)

## Results Analysis
After testing, performance comparisons across frameworks will be displayed. Compare:

- Requests per second (RPS)
- Average response time
- P99 response time (99% of requests respond within this time)
- Error rate and failed requests

## Extending Tests
You can extend tests by:

1. Adding new scenarios in internal/benchmark/scenario.go
2. Implementing corresponding route handlers in each framework
3. Adj3sting test parameters and analysis methods as needed

## Notes

- For fair comparison, all frameworks are configured with minimal middleware usage, disabling default logging and error handling middleware
- Benchmark results may be affected by hardware, network and system load
- In production environments, framework selection should also consider features, ecosystem and maintenance factors
