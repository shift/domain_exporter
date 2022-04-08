local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local promgrafonnet = import 'github.com/kubernetes-monitoring/kubernetes-mixin/lib/promgrafonnet/promgrafonnet.libsonnet';
local gauge = promgrafonnet.gauge;

{
  grafanaDashboards+:: {
    'domains.json':
      local domainExpiration =
        graphPanel.new(
          'Domains',
          datasource='$datasource',
          span=6,
          format='percentunit',
          max=1,
          min=0,
          stack=true,
        )
        .addTarget(prometheus.target(
          |||
            rate(domain_expiration[$__rate_interval])
          ||| % $._config,
          legendFormat='{{Domains}}',
          intervalFactor=5,
        ));

      dashboard.new(
        '%sDomains' % $._config.dashboardNamePrefix,
        time_from='now-1h',
        tags=($._config.dashboardTags),
        timezone='utc',
        refresh='30s',
        graphTooltip='shared_crosshair'
      )
      .addTemplate(
        {
          current: {
            text: 'default',
            value: 'default',
          },
          hide: 0,
          label: 'Data Source',
          name: 'datasource',
          options: [],
          query: 'prometheus',
          refresh: 1,
          regex: '',
          type: 'datasource',
        },
      )
      .addTemplate(
        template.new(
          'domain',
          '$datasource',
          'label_values(domain_expiration{%(domainExporterSelector)s}, domain)' % $._config,
          refresh='time',
        )
      )
      .addRow(
        row.new()
        .addPanel(domainExpiration)
      )
      .addRow(
        row.new()
        .addPanel(domainExpiration)
      ),
  },
}
