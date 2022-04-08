# Domain Mixin

The Domain Mixin is a set of configurable, reusable, and extensible alerts and
dashboards based on the metrics exported by the Domain Exporter. The mixin
creates alerting rules for Prometheus and suitable dashboard descriptions for
Grafana.

To use them, you need to have `jsonnet` (v0.16+) and `jb` installed. If you
have a working Go development environment, it's easiest to run the following:
```bash
$ go install github.com/google/go-jsonnet/cmd/jsonnet@latest
$ go install github.com/google/go-jsonnet/cmd/jsonnetfmt@latest
$ go install github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb@latest
```

Next, install the dependencies by running the following command in this
directory:
```bash
$ jb install
```

You can then build the Prometheus rules file `domain_alerts.yaml`:
```bash
$ make domain_alerts.yaml
```

You can also build a directory `dashboard_out` with the JSON dashboard files
for Grafana:
```bash
$ make dashboards_out
```

For more advanced uses of mixins, see
https://github.com/monitoring-mixins/docs.

