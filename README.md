# SIS Project - Checkpoint 1 (Text Analyzer)

## Project Overview

This is a simple synchronous, in-memory Go application that models a domain-specific task.
The chosen domain is a **Text Analyzer** that calculates the number of words and characters in a given string.

## How to Run

Since this is Checkpoint 1, there is no HTTP server or external interface yet.
To verify the application's correctness, run the automated tests:

```bash
go test -v ./...
```

To format and check the code, run:

gofmt -l .
go vet ./...

Task Lifecycle
The application tracks the state of each task. The allowed transitions are:

1. Pending - The task has been created with valid input but has not yet been processed.
2. Processing - The task is currently being processed (calculating words and characters).
3. Completed - The task has successfully completed processing, and the result is available.

Invalid transitions (e.g., trying to process an already "Completed" task) are blocked and will return an error.
