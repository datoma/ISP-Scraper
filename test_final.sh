#!/bin/bash

echo "=== Final ISP Scraper Test ==="

# Kill any existing processes
pkill -f isp-scraper 2>/dev/null || true

echo "Building scraper..."
go build -o isp-scraper main.go

echo "Starting scraper in background..."
./isp-scraper &
SCRAPER_PID=$!

echo "Waiting for scraper to initialize and collect data..."
sleep 8

echo ""
echo "=== Checking Metrics ==="
echo "Getting metrics from http://localhost:9090/metrics"

# Get specific ipinfo metrics
METRICS=$(curl -s http://localhost:9090/metrics | grep -E "(ipinfo_|# HELP)")

if [ -z "$METRICS" ]; then
    echo "ERROR: No ipinfo metrics found!"
    echo "Full metrics output:"
    curl -s http://localhost:9090/metrics | head -30
else
    echo "SUCCESS: Found ipinfo metrics:"
    echo "$METRICS"
fi

# Kill scraper
kill $SCRAPER_PID 2>/dev/null || true

echo ""
echo "=== Test Complete ==="