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