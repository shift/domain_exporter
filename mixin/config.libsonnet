{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.
    domainExporterSelector: 'job="domain"',

    domainRemainingDaysCriticalThreshold: 15,
    domainRemainingDaysWarningThreshold: 45,

    rateInterval: '5m',
    dashboardNamePrefix: 'Domain Exporter / ',
    dashboardTags: ['domain-exporter-mixin'],
  },
}
