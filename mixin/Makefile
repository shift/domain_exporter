JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 2 --string-style s --comment-style s

all: fmt domain_alerts.yaml dashboards_out lint

fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

domain_alerts.yaml: mixin.libsonnet config.libsonnet $(wildcard alerts/*)
	jsonnet -S alerts.jsonnet > $@

dashboards_out: mixin.libsonnet config.libsonnet $(wildcard dashboards/*)
	@mkdir -p dashboards_out
	jsonnet -J vendor -m dashboards_out dashboards.jsonnet

lint: domain_alerts.yaml
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" -; \
		done

	promtool check rules domain_alerts.yaml

.PHONY: jb_install
jb_install:
	jb install

clean:
	rm -rf dashboards_out domain_alerts.yaml
