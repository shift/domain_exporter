FROM scratch

MAINTAINER Vincent Palmer <@shift>

ADD docker/domain_exporter /

ENTRYPOINT ["/domain_exporter"]
