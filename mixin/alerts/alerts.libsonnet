{
  prometheusAlerts+:: {
    groups+: [
      {
        name: 'domain-exporter',
        rules: [
          {
            alert: 'DomainExpiringWarning',
            expr: |||
              domain_expiration{} < %(domainRemainingDaysWarningThreshold)d
            ||| % $._config,
            'for': '24h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Domain name expiration in less than %(domainRemainingDaysWarningThreshold)d days.' % $._config,
              description: 'The domain "{{ $labels.domain }}" has only {{ $value }} days left before expiration.',
            },
          },
          {
            alert: 'DomainExpiringCritical',
            expr: |||
              domain_expiration{} < %(domainRemainingDaysCriticalThreshold)d
            ||| % $._config,
            'for': '24h',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'Domain name expiration in less than %(domainRemainingDaysCriticalThreshold)d days.' % $._config,
              description: 'The domain "{{ $labels.domain }}" has only {{ $value }} days left before expiration.',
            },
          },

          {
            alert: 'DomainUnfindable',
            expr: 'domain_expiration_unfindable > 0',
            'for': '24h',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'Domain Exporter issue',
              description: 'Unable to find or parse expiration for "{{ $labels.domain }}"',
            },
          },
          {
            alert: 'DomainMetricsAbsent',
            expr: 'absent(domain_expiration) > 0',
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Domain Exporter issue',
              description: 'Metrics for Domain Exporter are absent',
            },
          },
          {
            alert: 'DomainExperationNegative',
            expr: 'domain_expiration{} < 0',
            'for': '24h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Domain Exporter issue',
              description: 'Negative value parsed for "{{ $labels.domain }}", The WHOIS response probably changed or doesn\'t parse correctly. Please create an issue with the response and we\'ll add the format.',
            },
          },
        ],
      },
    ],
  },
}
