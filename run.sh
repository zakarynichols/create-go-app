#!/bin/bash

# Navigate to the project directory (optional)
# cd /path/to/your/go/project

# Run the Go application with the specified arguments
go run . --type=http my-app

# Exit with the same status as the Go command
exit $?
