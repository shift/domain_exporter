# domain-exporter

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.1.12](https://img.shields.io/badge/AppVersion-v0.1.12-informational?style=flat-square)

domain_exporter chart

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| autoscaling | object | `{"enabled":false,"maxReplicas":100,"minReplicas":1,"targetCPUUtilizationPercentage":80}` | Leaving this in for the laughs |
| deployment.apiVersion | string | `""` | For those running 1.16 and less to define the apiVersion to be extensions/v1beta1, apps/v1beta1, or apps/v1beta2. Default: apps/v1 |
| domains[0] | string | `"google.com"` | List of domains to statically scrape. Used to populate ConfigMap used in the pod. |
| domains[1] | string | `"goolge.co.uk"` |  |
| domains[2] | string | `"google.de"` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/shift/domain_exporter"` |  |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.className | string | `""` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}]` | host ingress and path |
| ingress.tls | list | `[{"hosts":["chart-example.local"],"secretName":"chart-example-tls"}]` | tls secret ingress and hosts to use |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podMonitor.enabled | bool | `true` | If disabled prometheus.io/scrape will be set to true. |
| podSecurityContext | object | `{"fsGroup":2000}` | Try and be as secure as possible |
| replicaCount | int | `1` | Number of replicas to run when auto-scaling (ha) is turned off. |
| resources | object | `{"limits":{"cpu":"100m","memory":"32Mi"},"requests":{"cpu":"10m","memory":"8Mi"}}` | Define the limits and requests for cpu and memory, default values should be reasonable. |
| securityContext | object | `{"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsNonRoot":true,"runAsUser":1000}` | Try and be as secure as possible |
| service.port | int | `9203` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` |  |

