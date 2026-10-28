
# Go Channel Health Checker

A concurrent website health checker built with Go to explore
channels, goroutines, worker pools, and synchronization.

## Features

- Concurrent HTTP health checks
- Fixed-size worker pool
- Job distribution through channels
- Result collection through channels
- Send-only and receive-only channel parameters
- Graceful channel closure
- Worker synchronization using `sync.WaitGroup`
- HTTP request timeouts
- Race detection support

## Requirements

- Go 1.22+
- Internet access to reach the target websites

## Run

```bash
go mod init example.com/go-channel-health-checker
go run .
```

## Run with the race detector

```bash
go run -race .
```

## How It Works

1. The producer sends website URLs through the jobs channel.
2. Three workers receive jobs and perform HTTP requests concurrently.
3. Workers send their results through the results channel.
4. A WaitGroup tracks worker completion.
5. The coordinator closes the results channel after all workers finish.
6. The main goroutine consumes results until the channel closes.

## Concepts Demonstrated

### Unbuffered Channels

```go
jobs := make(chan string)
results := make(chan Result)
```

Sending and receiving synchronize when the channel is unbuffered.

### Channel Direction

```go
jobs <-chan string
results chan<- Result
```

Workers receive jobs and send results.

### Channel Closure

The producer closes `jobs` after sending all URLs.

The coordinator closes `results` only after every worker finishes,
preventing sends to a closed channel.

### WaitGroup

The WaitGroup coordinates completion of the workers. It does not
replace channels or protect arbitrary shared state.

## Important Notes

An HTTP status code of 404 is a valid HTTP response, but this
example classifies status codes of 400 or above as unhealthy.

The HTTP client timeout limits how long an individual request can
take. A production implementation may also need context-based
cancellation, retries, concurrency limits, and more detailed
error handling.

## Learning Goals

- Understand blocking and buffered channels.
- Build worker pools.
- Learn fan-out and fan-in patterns.
- Prevent deadlocks and goroutine leaks.
- Practice reliable concurrent programming in Go.
