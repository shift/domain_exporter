# domain_exporter

Very simple service which performs WHOIS lookups for a list of domains provided in the "config" file and exposes them on a "/metrics" endpoint for consumption via Prometheus.

````yaml
domains:
  - google.com
  - google.co.uk
````

Flags:
````bash
usage: domain_exporter [<flags>]

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --config="domains.yml"  Domain exporter configuration file.
      --bind=":9203"          The address to listen on for HTTP requests.
      --log.level=info        Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --version               Show application version. 
````

### Docker image

We publish a docker image [on the docker hub](https://hub.docker.com/r/shift/domain_exporter/). You can pull this with `docker pull shift/domain_exporter`.
