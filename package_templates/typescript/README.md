# Grafana Foundation SDK – TypeScript

A set of tools, types and *builder libraries* for building and manipulating Grafana objects in TypeScript.

> [!NOTE]
> This branch contains **types and builders generated for Grafana {{ .Extra.GrafanaVersion }}.**
> Other supported versions of Grafana can be found at [this repository's root](https://github.com/grafana/grafana-foundation-sdk/).

## Installing

```shell
yarn add '@grafana/grafana-foundation-sdk@~{{ .Extra.GrafanaVersion|registryToSemver }}-cog{{ .Extra.CogVersion }}.{{ .Extra.BuildTimestamp }}'
```

## Example usage

```typescript
import { DashboardBuilder, RowBuilder } from '@grafana/grafana-foundation-sdk/dashboard';
import { DataqueryBuilder } from '@grafana/grafana-foundation-sdk/prometheus';
import { PanelBuilder } from '@grafana/grafana-foundation-sdk/timeseries';

const builder = new DashboardBuilder('[TEST] Node Exporter / Raspberry')
  .uid('test-dashboard-raspberry')
  .tags(['generated', 'raspberrypi-node-integration'])

  .refresh('1m')
  .time({from: 'now-30m', to: 'now'})
  .timezone('browser')

  .withRow(new RowBuilder('Overview'))
  .withPanel(
    new PanelBuilder()
      .title('Network Received')
      .unit('bps')
      .min(0)
      .withTarget(
        new DataqueryBuilder()
          .expr('rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", device!="lo"}[$__rate_interval]) * 8')
          .legendFormat({{ `"{{ device }}"` }})
      )
  )
;

console.log(JSON.stringify(builder.build(), null, 2));
```

## Maturity

> [!WARNING]
> The code in this repository should be considered experimental. Documentation is only
available alongside the code. It comes with no support, but we are keen to receive
feedback and suggestions on how to improve it, though we cannot commit
to resolution of any particular issue.

Grafana Labs defines experimental features as follows:

> Projects and features in the Experimental stage are supported only by the Engineering
teams; on-call support is not available. Documentation is either limited or not provided
outside of code comments. No SLA is provided.
>
> Experimental projects or features are primarily intended for open source engineers who
want to participate in ensuring systems stability, and to gain consensus and approval
for open source governance projects.
>
> Projects and features in the Experimental phase are not meant to be used in production
environments, and the risks are unknown/high.

## License

[Apache 2.0 License](./LICENSE)
