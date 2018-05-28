# domain_exporter

Very simple service which performs WHOIS lookups for a list of domains provided in the "config" file and exposes them on a "/metrics" endpoint for consumption via Prometheus.

```yaml
domains:
  - google.com
  - google.co.uk
```

Flags:
```bash
usage: domain_exporter [<flags>]

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --config="domains.yml"  Domain exporter configuration file.
      --bind=":9203"          The address to listen on for HTTP requests.
      --log.level=info        Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --version               Show application version.
```

### Docker image

We publish a docker image [on the Quay registry](https://quay.io/repository/shift/domain_exporter). You can pull this with `docker pull quay.io/shift/domain_exporter`.

### Example Prometheus Alert

The following alert will be triggered when domains expire within 45 days

```yaml
groups:
 - name: ./domain.rules
   rules:
    - alert: DomainExpiring
      expr: domain_expiration{} < 45
      for: 24h
      labels:
        severity: warning
      annotations:
        description: "{{ $labels.domain }} expires in {{ $value }} days"
```
