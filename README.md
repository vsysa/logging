# Go Logger Module

## Overview

This library provides a flexible logging interface along with several implementations, typically as wrappers around popular logging libraries.
The logger is designed to integrate easily into your applications, offering key features that enhance logging management and flexibility.

## Status

**Early Development Stage**: This module is currently in the early stages of development. As such, features and the API are subject to change.


## Key Features

- **Logger Interface** : A unified interface for logging, allowing seamless integration with different logging libraries.

- **Multiple Implementations** : Choose from several logger implementations, all adhering to the same interface, simplifying switching between different logging libraries.

- **Logging Context** : Supports logging with context, ensuring that important metadata and contextual information are preserved across logs.

- **Dynamic Log Level Setting** : Change log levels dynamically at runtime, providing control over the verbosity of logs without needing to restart the application.

- **Cloning Capabilities** : Easily create new logger instances that share the same context and configurations, streamlining the process of managing logger instances across different parts of your application.


## Installation

To include this module in your Go project, use the following `go get` command:

```bash
go get https://github.com/vsysa/logging@latest
```

### Implementations

- **Logrus** : A wrapper around the popular Logrus library.
- **Zap** : An implementation that wraps the high-performance Zap logger.
- **Test** : An implementation that can collect all logs and print them if necessary.

Logger implementations are located in `logging/logger/<logger implementation>`


For more information about configuration and usage, please refer to the documentation for each specific logger.


## Usage Examples

Below are various examples illustrating how to use the Go Logger module in your applications.

### Basic Logger Setup

```go
package main

func main() {
    // Initialize a new logger instance
    logger := logging.NewLogrusLogger()

    // Set the log level to Info
	logger.SetLevel(logging.InfoLevel)

    // Log a simple message
	logger.Info("This is an informational message.")
}
```



### Adding Context to Logs

```go
func main() {
    logger := logging.NewLogrusLogger()

    // Add single context
    logger.AddContext("user_id", "12345")

    // Add multiple contexts at once
    logger.AddContexts(logging.ContextParams{
        "order_id": "abcde",
        "payment_mode": "credit_card",
    })

    // Log a message with context
    logger.Info("Order processed successfully.")
}
```



### Working with Context and Cloning

```go
func main() {
    // Create a background context with a logger
    ctx := logctx.CtxWithLogger(context.Background(), logging.NewLogrusLogger())

    // Retrieve a logger from context, adding context parameters
    logger := logctx.L(ctx).AddContext("request_id", "xyz789")

    // Log a message using the logger from context
    logger.Info("Starting request handling")

    // Clone the logger from context and modify the clone
    clonedLogger := logging.LoggerCloneFromCtx(ctx).AddContext("cloned", true)
    clonedLogger.Info("This is from the cloned logger")

    // The original logger remains unaffected
    logger.Info("Finishing request handling")
}
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Authors
- **Vladislav Sysalov** - [vsysa](https://github.com/vsysa)