package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shift/whois"
	"gopkg.in/yaml.v2"
)

var (
	// How often to check domains
	checkRate = 12 * time.Hour

	httpBind         string
	configFile       string

	domainExpiration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "domain_expiration",
			Help: "Days until the WHOIS record states this domain will expire",
		},
		[]string{"domain"},
	)

	expiryRegex = regexp.MustCompile(`(?i)(Registry Expiry Date|paid-till|Expiration Date|Expiry.*|expires.*): (.*)`)

	formats = []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"02-Jan-2006",
		"2006.01.02",
		"Mon Jan 2 15:04:05 MST 2006",
		"02/01/2006",
		"2006-01-02 15:04:05 MST",
		"2006/01/02",
		"Mon Jan 2006 15:04:05",
		"2006-01-02 15:04:05-07",
	}
)

type Config struct {
	Domains []string `yaml:"domains"`
}

func main() {
	flag.StringVar(&httpBind, "bind", ":9203", "Port to expose the /metrics endpoint on")
	flag.StringVar(&configFile, "config", "./config", "Path to configuration file")
	flag.Parse()
	prometheus.Register(domainExpiration)

	config := Config{}

	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			for _, query := range config.Domains {
				req, err := whois.NewRequest(query)
				if err != nil {
					log.Print(err)
				}

				var res *whois.Response
				res, err = whois.DefaultClient.Fetch(req)
				if err != nil {
					log.Print(err)
					continue

				}

				result := expiryRegex.FindStringSubmatch(res.String())

				if len(result) < 2 {
					log.Print(fmt.Sprintf("Don't know how to parse domain: %s\n", query))
					continue
				}

				parsed := false
				for _, format := range formats {
					if date, err := time.Parse(format, strings.TrimSpace(result[2])); err == nil {
						days := math.Floor(date.Sub(time.Now()).Hours() / 24)
						log.Print(fmt.Sprintf("Domain: %s, Days: %v, Date: %s", query, days, date))
						domainExpiration.WithLabelValues(query).Set(days)
						parsed = true
						break
					}

				}
				if !parsed {
					log.Print(fmt.Sprintf("Unable to parse date: %s, for %s\n", strings.TrimSpace(result[2]), query))
				}
			}
			time.Sleep(12 * time.Hour)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s\n", httpBind)
	log.Fatal(http.ListenAndServe(httpBind, nil))
}
