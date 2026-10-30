#!/bin/bash

# Simple script to run the ISP scraper
echo "Starting ISP Scraper..."
echo "Make sure you have set IPINFO_TOKEN environment variable"
echo "Example: export IPINFO_TOKEN=your_token_here"

if [ -z "$IPINFO_TOKEN" ]; then
    echo "Error: IPINFO_TOKEN is not set"
    exit 1
fi

go run main.go