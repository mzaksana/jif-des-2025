#!/bin/bash

# Stress test script for scaling demo
# Usage: ./stress-test.sh [requests] [concurrency]

REQUESTS=${1:-1000}
CONCURRENCY=${2:-50}
URL="http://localhost:8080/api/work"

echo "========================================"
echo "  Stress Test Configuration"
echo "========================================"
echo "  URL:         $URL"
echo "  Requests:    $REQUESTS"
echo "  Concurrency: $CONCURRENCY"
echo "========================================"
echo ""

# Check if hey is installed
if command -v hey &> /dev/null; then
    echo "Using 'hey' for load testing..."
    echo ""
    hey -n $REQUESTS -c $CONCURRENCY $URL
elif command -v wrk &> /dev/null; then
    echo "Using 'wrk' for load testing..."
    echo ""
    wrk -t4 -c$CONCURRENCY -d30s $URL
elif command -v ab &> /dev/null; then
    echo "Using 'ab' (Apache Bench) for load testing..."
    echo ""
    ab -n $REQUESTS -c $CONCURRENCY $URL
else
    echo "No load testing tool found!"
    echo ""
    echo "Please install one of the following:"
    echo "  - hey:  go install github.com/rakyll/hey@latest"
    echo "  - wrk:  brew install wrk"
    echo "  - ab:   comes with Apache (apache2-utils)"
    exit 1
fi
