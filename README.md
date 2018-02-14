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
        Port to expose the /metrics endpoint on (default ":9203")
  -config string
        Path to configuration file (default "./config")

````

### Docker image

We publish a docker image [on the docker hub](https://hub.docker.com/r/shift/domain_exporter/). You can pull this with `docker pull shift/domain_exporter`.
