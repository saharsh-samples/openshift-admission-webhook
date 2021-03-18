#!/bin/bash

# Define script constants
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
COVERAGE_OUT="$SCRIPT_DIR/coverage.out"
COVERAGE_HTML="$SCRIPT_DIR/coverage.html"

# Process CLI args
CREATE_HTML_REPORT=0
if [ "$1" == "--html" ]; then
    CREATE_HTML_REPORT=1
fi

# Change to script directory and turn off exit on error
pushd $SCRIPT_DIR
set +e

# Run tests and capture pass/fail
go test -coverprofile "$COVERAGE_OUT" ./...
TESTS_PASSED=$?
go tool cover -func "$COVERAGE_OUT" | grep "total:"

# If configured, create HTML report of coverage
if [[ $TESTS_PASSED -eq 0 && $CREATE_HTML_REPORT -eq 1 ]]; then
    echo "Creating HTML Report.."
    go tool cover -html="$COVERAGE_OUT" -o "$COVERAGE_HTML"
fi

# Turn exit on error back on and restore original directory
set -e
popd
exit $TESTS_PASSED