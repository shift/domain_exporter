FROM scratch
COPY domain_exporter /domain_exporter
ENTRYPOINT ["/domain_exporter"]
