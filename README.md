# domain_exporter

Very simple service which performs WHOIS lookups for a list of domains provided in the "config" file and exposes them on a "/metrics" endpoint for consumption via Prometheus.

````yaml
domains:
  - google.com
  - google.co.uk
````

Flags:
````bash
Usage of ./domain_exporter:
  -bind string
        Port to expose the /metrics endpoint on (default ":9080")
  -config string
        Path to configuration file (default "./config")

````
