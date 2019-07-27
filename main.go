package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/domainr/whois"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
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
	unfindableExpiration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "domain_expiration_unfindable",
			Help: "That the domain date could not be parsed, or the domain doesn't have a whois record",
		},
		[]string{"domain"},
	)

	expiryRegex = regexp.MustCompile(`(?i)(Registry Expiry Date|paid-till|Expiration Date|Expiration Time|Expiry.*|expires.*): (.*)`)

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
		"2006-01-02 15:04:05",
		"2.1.2006 15:04:05",
	}

	config promlog.Config
	logger log.Logger
)

type Config struct {
	Domains []string `yaml:"domains"`
}

func main() {
	flag.AddFlags(kingpin.CommandLine, &config)
	kingpin.Version(version.Print("domain_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger = promlog.New(&config)

	level.Info(logger).Log("msg", "Starting domain_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	prometheus.Register(domainExpiration)
	prometheus.Register(unfindableExpiration)

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
					_, err = lookup(query, domainExpiration, unfindableExpiration)
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
	probeUnfindableExpiration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "domain_expiration_unfindable",
			Help: "That the domain date could not be parsed, or the domain doesn't have a whois record",
		},
		[]string{"domain"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeExpiration)
	registry.MustRegister(probeUnfindableExpiration)
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}
	_, err := lookup(target, probeExpiration, unfindableExpiration)
	if err != nil {
		level.Warn(logger).Log("warn", err)
		http.Error(w, fmt.Sprintf("Don't know how to parse: %q", target), http.StatusBadRequest)
		return
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func parse(host string, res []byte) (float64, error) {
	results := expiryRegex.FindStringSubmatch(string(res))
	if len(results) < 1 {
		err := fmt.Errorf("Don't know how to parse domain: %s", host)
		level.Warn(logger).Log("warn", err.Error())
		return -2, err
	}

	for _, format := range formats {
		if date, err := time.Parse(format, strings.TrimSpace(results[2])); err == nil {
			days := math.Floor(date.Sub(time.Now()).Hours() / 24)
			level.Info(logger).Log("domain:", host, "days", days, "date", date)
			return days, nil
		}

	}
	return -1, errors.New(fmt.Sprintf("Unable to parse date: %s, for %s\n", strings.TrimSpace(results[2]), host))
}

func lookup(domain string, handler *prometheus.GaugeVec, badhandler *prometheus.GaugeVec) (float64, error) {
	req, err := whois.NewRequest(domain)
	if err != nil {
		return -1, err
	}

	res, err := whois.DefaultClient.Fetch(req)
	if err != nil {
		return -1, err
	}

	days, err := parse(domain, res.Body)
	if err != nil {
		if badhandler != nil {
			badhandler.WithLabelValues(domain).Set(1)
		}
		return -1, err
	}

	if handler != nil {
		handler.WithLabelValues(domain).Set(days)
	}
	if badhandler != nil {
		badhandler.WithLabelValues(domain).Set(0)
	}

	return days, nil
}
