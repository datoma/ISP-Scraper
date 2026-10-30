package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tidwall/gjson"
)

var (
	ipInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ipinfo_data",
		Help: "IP information from ipinfo.io API",
	},
		[]string{"ip", "asn", "as_name", "as_domain", "country_code", "country", "continent_code", "continent"},
	)

	ipInfoTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ipinfo_timestamp",
		Help: "Timestamp of the scrape",
	})

	lastScrapeTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ipinfo_last_scrape_timestamp",
		Help: "Timestamp of the last successful scrape",
	})

	scrapeSuccess = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ipinfo_scrape_success",
		Help: "1 if the last scrape was successful, 0 otherwise",
	})
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := os.Getenv("IPINFO_TOKEN")
	if token == "" {
		log.Fatal("IPINFO_TOKEN environment variable is required")
	}

	// Get scrape interval from environment variable (default 600 seconds = 10 minutes)
	scrapeInterval := 600
	if intervalStr := os.Getenv("SCRAPE_INTERVAL_SECONDS"); intervalStr != "" {
		var interval int
		if _, err := fmt.Sscanf(intervalStr, "%d", &interval); err == nil && interval > 0 {
			scrapeInterval = interval
		} else {
			log.Println("Invalid SCRAPE_INTERVAL_SECONDS, using default 600 seconds")
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		scrapeData(token)
		time.Sleep(time.Duration(scrapeInterval) * time.Second)
	}
}

func scrapeData(token string) {
	url := "https://api.ipinfo.io/lite/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		scrapeSuccess.Set(0)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		scrapeSuccess.Set(0)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status code %d", resp.StatusCode)
		scrapeSuccess.Set(0)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		scrapeSuccess.Set(0)
		return
	}

	jsonData := string(body)

	ip := gjson.Get(jsonData, "ip").String()
	asn := gjson.Get(jsonData, "asn").String()
	asName := gjson.Get(jsonData, "as_name").String()
	asDomain := gjson.Get(jsonData, "as_domain").String()
	countryCode := gjson.Get(jsonData, "country_code").String()
	country := gjson.Get(jsonData, "country").String()
	continentCode := gjson.Get(jsonData, "continent_code").String()
	continent := gjson.Get(jsonData, "continent").String()

	ipInfo.WithLabelValues(
		ip,
		asn,
		asName,
		asDomain,
		countryCode,
		country,
		continentCode,
		continent,
	).Set(1)

	ipInfoTimestamp.Set(float64(time.Now().Unix()))

	lastScrapeTimestamp.Set(float64(time.Now().Unix()))
	scrapeSuccess.Set(1)

	log.Printf("Scraped data for IP: %s, ASN: %s, AS Name: %s", ip, asn, asName)
}
