# ISP Scraper

A Prometheus scraper that fetches ISP information from ipinfo.io and exposes it as metrics.

## Features

- Fetches ISP data from ipinfo.io API
- Exposes each data field as a separate Prometheus metric:
- Configurable scrape interval via environment variable


## Usage

1. Create a `.env` file with your ipinfo.io API token:
   ```bash
   echo "IPINFO_TOKEN=your_token_here" > .env
   ```

2. Set scrape interval (optional, defaults to 10 minutes):
   ```bash
   echo "SCRAPE_INTERVAL_SECONDS=600" >> .env
   ```

3. Run the scraper:
   ```bash
   go run main.go
   ```

4. Access metrics at http://localhost:9090/metrics

## Metrics Exposed

- `ipinfo_ip` - IP address from ipinfo.io API
- `ipinfo_asn` - Autonomous System Number from ipinfo.io API  
- `ipinfo_as_name` - Autonomous System name from ipinfo.io API
- `ipinfo_as_domain` - Autonomous System domain from ipinfo.io API
- `ipinfo_country_code` - Country code from ipinfo.io API
- `ipinfo_country` - Country name from ipinfo.io API
- `ipinfo_continent_code` - Continent code from ipinfo.io API
- `ipinfo_continent` - Continent name from ipinfo.io API
- `ipinfo_timestamp` - Timestamp of the scrape
- `ipinfo_last_scrape_timestamp` - Unix timestamp of last successful scrape
- `ipinfo_scrape_success` - 1 if last scrape was successful, 0 otherwise

## Configuration

- `IPINFO_TOKEN` (required): Your ipinfo.io API token
- `SCRAPE_INTERVAL_SECONDS` (optional): How often to scrape in seconds (default: 600 = 10 minutes)

## Docker Deployment

```bash
docker build -t isp-scraper .
docker run -v $(pwd)/.env:/app/.env -p 9090:9090 isp-scraper
```

## Grafana Dashboard

A sample Grafana dashboard is included in the `grafana-dashboard.json` file that visualizes ISP usage over time.