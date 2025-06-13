# domain_exporter
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fshift%2Fdomain_exporter.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fshift%2Fdomain_exporter?ref=badge_shield)


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
      --log.format=logfmt     Output format of log messages. One of: [logfmt, json]
      --version               Show application version.
```

### Docker image

We publish a docker image [on the Quay registry](https://quay.io/repository/shift/domain_exporter). You can pull this with `docker pull ghcr.io/shift/domain_exporter`.

### Running on Kubernetes

[Here](contrib/k8s-domain-exporter.yaml) is an example Kubernetes deployment configuration for how to deploy the domain_exporter.

### Probe Endpoint

In addition to the `/metrics` endpoint which exposes metrics for the configured domains, the exporter also provides a `/probe` endpoint which allows on-demand querying of specific domains.

Example usage:
```
http://localhost:9203/probe?target=google.com
```

### Prometheus Configuration

You can configure Prometheus to use the domain_exporter with the following configuration:

```yaml
- job_name: 'domain-exporter'
  metrics_path: /probe
  static_configs:
    - targets:
      - google.com  # Target to probe with whois.
      - facebook.com  # Target to probe with whois.
  relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: localhost:9203  # The domain exporter's real hostname:port.
```

### Example Prometheus Alert

The following alert will be triggered when domains expire within 45 days, or if
they don't have a whois record available (perhaps having been long expired).

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
    - alert: DomainUnfindable
      expr: domain_expiration_unfindable > 0
      for: 24h
      labels:
        severity: critical
      annotations:
        description: "Unable to find or parse expiry for {{ $labels.domain }}"
    - alert: DomainMetricsAbsent
      expr: absent(domain_expiration) > 0
      for: 1h
      labels:
        severity: warning
      annotations:
        description: "Metrics for domain-exporter are absent"
```

### FAQ

##### Why did I get a negative amount of days until expiry?

The WHOIS resposne probably doesn't parse correctly. Please create an issue with the response and we'll add the format.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fshift%2Fdomain_exporter.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fshift%2Fdomain_exporter?ref=badge_large)
