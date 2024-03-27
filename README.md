# Go Logger Module

## Overview

This module provides a flexible and extensible logging solution for Go applications, wrapping the popular `logrus` library to offer enhanced functionality. Designed with simplicity and efficiency in mind, it supports context management, dynamic log level setting, timer functions for performance monitoring, and cloning capabilities to create logger instances with shared contexts and configurations.

### Status

**Early Development Stage**: This module is currently in the early stages of development. As such, features and the API are subject to change.

### Features

- **Context Management**: Easily add, modify, and delete logging context to enrich log messages.
- **Dynamic Log Levels**: Adjust log verbosity on the fly, tailoring output to your current needs.
- **Timer Functions**: Integrated timer functions to measure and log operation durations, aiding in performance analysis.
- **Cloning**: Create clones of your logger instances, preserving context and settings while facilitating isolated modifications.

## Installation

To include this module in your Go project, use the following `go get` command:

```bash
go get https://github.com/vsysa/logging@latest
```

## Roadmap
**Open Telemetry Integration**: Plan to integrate with OpenTelemetry to support tracing and metrics, providing a comprehensive observability solution.
Contributing


## Usage Examples

Below are various examples illustrating how to use the Go Logger module in your applications.

### Basic Logger Setup

```go
package main

import (
    "github.com/vsysa/logging"
)

func main() {
    // Initialize a new logger instance
    log := logging.NewBaseLogger()

    // Set the log level to Info
    log.SetLevel(logging.InfoLevel)

    // Log a simple message
    log.Info("This is an informational message.")
}
```



### Adding Context to Logs

```go
func main() {
    log := logging.NewBaseLogger()

    // Add single context
    log.AddContext("user_id", "12345")

    // Add multiple contexts at once
    log.AddContexts(logging.ContextParams{
        "order_id": "abcde",
        "payment_mode": "credit_card",
    })

    // Log a message with context
    log.Info("Order processed successfully.")
}
```



### Using Timers for Performance Monitoring

```go
func main() {
    log := logging.NewBaseLogger()

    // Start a timer
    log.TimerStart("db_query")

    // Simulate some operation, e.g., database query
    time.Sleep(2 * time.Second)

    // Print the elapsed time directly
    log.TimerPrint("db_query")

    // Or, get the duration for manual handling
    duration, ok := log.TimerDuration("db_query")
    if ok {
        log.Info(fmt.Sprintf("DB query took %v milliseconds", duration))
    }
}
```



### Working with Context and Cloning

```go
func main() {
    // Create a background context with a logger
    ctx := logging.BackgroundCtxLogger(logging.NewBaseLogger())

    // Retrieve a logger from context, adding context parameters
    logger := logging.L(ctx).AddContext("request_id", "xyz789")

    // Log a message using the logger from context
    logger.Info("Starting request handling")

    // Clone the logger from context and modify the clone
    clonedLogger := logging.LoggerCloneFromCtx(ctx).AddContext("cloned", true)
    clonedLogger.Info("This is from the cloned logger")

    // The original logger remains unaffected
    logger.Info("Finishing request handling")
}
```