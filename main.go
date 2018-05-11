package main

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/shift/whois"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	// How often to check domains
	checkRate = 12 * time.Hour

	configFile = kingpin.Flag("config", "Domain exporter configuration file.").Default("domains.yml").String()
	httpBind   = kingpin.Flag("bind", "The address to listen on for HTTP requests.").Default(":9203").String()

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
	allowedLevel := promlog.AllowedLevel{}
	flag.AddFlags(kingpin.CommandLine, &allowedLevel)
	kingpin.Version(version.Print("domain_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(allowedLevel)
	level.Info(logger).Log("msg", "Starting domain_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())
	prometheus.Register(domainExpiration)

	config := Config{}

	filename, err := filepath.Abs(*configFile)
	if err != nil {
		level.Warn(logger).Log("warn", err)
	}
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		level.Warn(logger).Log("warn", err)
		level.Warn(logger).Log("warn", "Configuration file not present, you'll have to /probe me for metrics.")
	}
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		level.Warn(logger).Log("warn", err)
	} else {
		go func() {
			for {
				for _, query := range config.Domains {
					_, err = lookup(query, domainExpiration, logger)
					if err != nil {
						level.Warn(logger).Log("warn", err)
					}
					continue
				}
				time.Sleep(checkRate)
			}
		}()
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, logger)
	})
	level.Info(logger).Log("msg", "Listening", "port", *httpBind)
	if err := http.ListenAndServe(*httpBind, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func probeHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	probeExpiration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "domain_expiration",
			Help: "Days until the WHOIS record states this domain will expire",
		},
		[]string{"domain"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeExpiration)
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}
	_, err := lookup(target, probeExpiration, logger)
	if err != nil {
		level.Warn(logger).Log("warn", err)
		http.Error(w, fmt.Sprintf("Don't know how to parse: %q", target), http.StatusBadRequest)
		return
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func lookup(domain string, handler *prometheus.GaugeVec, logger log.Logger) (float64, error) {
	req, err := whois.NewRequest(domain)
	if err != nil {
		return -1, err
	}

	var res *whois.Response
	res, err = whois.DefaultClient.Fetch(req)
	if err != nil {
		return -1, err
	}

	result := expiryRegex.FindStringSubmatch(res.String())

	if len(result) < 2 {
		level.Warn(logger).Log("warn", fmt.Sprintf("Don't know how to parse domain: %s\n", domain))
		return -1, nil
	}

	for _, format := range formats {
		if date, err := time.Parse(format, strings.TrimSpace(result[2])); err == nil {
			days := math.Floor(date.Sub(time.Now()).Hours() / 24)
			level.Info(logger).Log("domain:", domain, "days", days, "date", date)
			handler.WithLabelValues(domain).Set(days)
			return days, nil
		}

	}
	return -1, errors.New(fmt.Sprintf("Unable to parse date: %s, for %s\n", strings.TrimSpace(result[2]), domain))
}
