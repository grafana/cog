# Grafana Foundation SDK – PHP

A set of tools, types and *builder libraries* for building and manipulating Grafana objects in PHP.

> [!NOTE]
> This branch contains **types and builders generated for Grafana {{ .Extra.GrafanaVersion }}.**
> Other supported versions of Grafana can be found at [this repository's root](https://github.com/grafana/grafana-foundation-sdk/).

## Installing

```shell
composer require "grafana/foundation-sdk:{{ .Extra.ReleaseBranch }}"
```

## Example usage

[More examples](https://github.com/grafana/grafana-foundation-sdk/tree/main/examples/php) can be found at the repository root.

### Building a dashboard

```php
<?php

use Grafana\Foundation\Common;
use Grafana\Foundation\Dashboard\DashboardBuilder;
use Grafana\Foundation\Dashboard\RowBuilder;
use Grafana\Foundation\Prometheus;
use Grafana\Foundation\Timeseries;

require_once __DIR__.'/vendor/autoload.php';

$builder = (new DashboardBuilder(title: 'Sample dashboard'))
    ->uid('generated-from-php')
    ->tags(['generated', 'from', 'php'])
    ->refresh('1m')
    ->time('now-30m', 'now')
    ->timezone(Common\Constants::TIME_ZONE_BROWSER)
    ->withRow(new RowBuilder('Overview'))
    ->withPanel(
        (new Timeseries\PanelBuilder())
            ->title('Network received')
            ->unit('bps')
            ->min(0)
            ->withTarget(
                (new Prometheus\DataqueryBuilder())
                    ->expr('rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", device!="lo"}[$__rate_interval]) * 8')
                    ->legendFormat({{ `'{{ device }}'` }})
            )
    )
;

echo(json_encode($builder->build(), JSON_PRETTY_PRINT).PHP_EOL);
```

### Unmarshaling a dashboard

```php
<?php

use Grafana\Foundation\Dashboard\Dashboard;

require_once __DIR__.'/vendor/autoload.php';

$dashboardJSON = file_get_contents(__DIR__.'/dashboard.json');

$dashboard = Dashboard::fromArray(json_decode($dashboardJSON, true));

var_dump($dashboard);

```

### Defining a custom query type

While the SDK ships with support for all core datasources and their query types,
it can be extended for private/third-party plugins.

To do so, define a type for the custom query:

```php
<?php

use Grafana\Foundation\Cog;

class CustomQuery implements \JsonSerializable, Cog\Dataquery
{
    // RefId and Hide are expected on all queries
    public string $refId;
    public ?bool $hide;

    // Query is specific to the CustomQuery type
    public string $expr;

    /**
     * @param string|null $expr
     * @param string|null $refId
     * @param bool|null $hide
     */
    public function __construct(?string $expr = null, ?string $refId = null, ?bool $hide = null)
    {
        $this->expr = $expr ?: "";
        $this->refId = $refId ?: "";
        $this->hide = $hide;
    }

    /**
     * @param array{expr?: string, refId?: string, hide?: bool} $inputData
     */
    public static function fromArray(array $data): self
    {
        return new self(
            expr: $data["expr"] ?? null,
            refId: $data["refId"] ?? null,
            hide: $data["hide"] ?? null,
        );
    }

    public function jsonSerialize(): array
    {
        $data = [
            "expr" => $this->expr,
            "refId" => $this->refId,
        ];
        if (isset($this->hide)) {
            $data["hide"] = $this->hide;
        }
        return $data;
    }
}
```

Now, let's define a builder for that type:

```php
<?php

use Grafana\Foundation\Cog;

/**
 * @implements Cog\Builder<CustomQuery>
 */
class CustomQueryBuilder implements Cog\Builder
{
    protected CustomQuery $internal;

    public function __construct(string $query)
    {
    	$this->internal = new CustomQuery(expr: $query);
    }

    /**
     * @return CustomQuery
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * The actual expression/query that will be evaluated by Prometheus
     */
    public function expr(string $expr): static
    {
        $this->internal->expr = $expr;
        return $this;
    }

    /**
     * A unique identifier for the query within the list of targets.
     * In server side expressions, the refId is used as a variable name to identify results.
     * By default, the UI will assign A->Z; however setting meaningful names may be useful.
     */
    public function refId(string $refId): static
    {
        $this->internal->refId = $refId;
        return $this;
    }

    /**
     * If hide is set to true, Grafana will filter out the response(s) associated with this query before returning it to the panel.
     */
    public function hide(bool $hide): static
    {
        $this->internal->hide = $hide;
        return $this;
    }
}
```

Register the type with cog, and use it as usual to build a dashboard:

```php
<?php

use Grafana\Foundation\Cog;
use Grafana\Foundation\Dashboard\DashboardBuilder;
use Grafana\Foundation\Dashboard\RowBuilder;
use Grafana\Foundation\Timeseries;

require_once __DIR__.'/vendor/autoload.php';

// This lets cog know about the newly created query type and how to unmarshal it.
Cog\Runtime::get()->registerDataqueryVariant(new Cog\DataqueryConfig(
    identifier: 'custom', // datasource plugin ID
    fromArray: [CustomQuery::class, 'fromArray'],
));

$builder = (new DashboardBuilder(title: 'Custom query type'))
    ->uid('test-custom-query-type')
    ->refresh('1m')
    ->time('now-30m', 'now')
    ->withRow(new RowBuilder('Overview'))
    ->withPanel(
        (new Timeseries\PanelBuilder())
            ->title('Sample panel')
            ->withTarget(new CustomQueryBuilder('query here'))
    )
;

echo(json_encode($builder->build(), JSON_PRETTY_PRINT).PHP_EOL);
```

### Defining a custom panel type

While the SDK ships with support for all core panels, it can be extended for
private/third-party plugins.

To do so, define a type and a builder for the custom panel's options:

```go
package main

import (
    "encoding/json"

    "github.com/grafana/grafana-foundation-sdk/go/cog"
    cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
    "github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

type CustomPanelOptions struct {
    MakeBeautiful bool `json:"makeBeautiful"`
}

func CustomPanelVariantConfig() cogvariants.PanelcfgConfig {
    return cogvariants.PanelcfgConfig{
        Identifier: "custom-panel", // plugin ID
        OptionsUnmarshaler: func(raw []byte) (any, error) {
            options := CustomPanelOptions{}

            if err := json.Unmarshal(raw, &options); err != nil {
                return nil, err
            }

            return options, nil
        },
    }
}

// Compile-time check to ensure that CustomPanelBuilder indeed is
// a builder for a dashboard.Panel.
var _ cog.Builder[dashboard.Panel] = (*CustomPanelBuilder)(nil)

type CustomPanelBuilder struct {
    internal *dashboard.Panel
    errors   map[string]cog.BuildErrors
}

func NewCustomPanelBuilder() *CustomPanelBuilder {
    return &CustomPanelBuilder{
        internal: &dashboard.Panel{
            Type: "custom-panel",
        },
        errors: make(map[string]cog.BuildErrors),
    }
}

func (builder *CustomPanelBuilder) Build() (dashboard.Panel, error) {
    var errs cog.BuildErrors

    for _, err := range builder.errors {
        errs = append(errs, cog.MakeBuildErrors("CustomPanel", err)...)
    }

    if len(errs) != 0 {
        return dashboard.Panel{}, errs
    }

    return *builder.internal, nil
}

func (builder *CustomPanelBuilder) Title(title string) *CustomPanelBuilder {
    builder.internal.Title = &title

    return builder
}

func (builder *CustomPanelBuilder) WithTarget(targets cog.Builder[cogvariants.Dataquery]) *CustomPanelBuilder {
    targetsResource, err := targets.Build()
    if err != nil {
        builder.errors["targets"] = err.(cog.BuildErrors)
        return builder
    }
    builder.internal.Targets = append(builder.internal.Targets, targetsResource)

    return builder
}

// [other panel options omitted for brevity]

func (builder *CustomPanelBuilder) MakeBeautiful() *CustomPanelBuilder {
    if builder.internal.Options == nil {
        builder.internal.Options = &CustomPanelOptions{}
    }
    builder.internal.Options.(*CustomPanelOptions).MakeBeautiful = true

    return builder
}
```

Register the type with cog, and use it as usual to build a dashboard:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/grafana/grafana-foundation-sdk/go/cog"
    "github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
    "github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func main() {
    // Required to correctly unmarshal panels and dataqueries
    plugins.RegisterDefaultPlugins()

    // This lets cog know about the newly created panel type and how to unmarshal it.
    cog.NewRuntime().RegisterPanelcfgVariant(CustomPanelVariantConfig())

    sampleDashboard, err := dashboard.NewDashboardBuilder("Custom panel type").
        Uid("test-custom-panel").
        Refresh("1m").
        Time("now-30m", "now").
        WithRow(dashboard.NewRowBuilder("Overview")).
        WithPanel(
            NewCustomPanelBuilder().
                Title("Sample panel").
                MakeBeautiful(),
        ).
        Build()
    if err != nil {
        panic(err)
    }

    dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
    if err != nil {
        panic(err)
    }

    fmt.Println(string(dashboardJson))
}
```

## Maturity

The code in this repository should be considered as "public preview" and is actively developed and maintained by Engineering teams at Grafana.

While this repository is stable enough to be used in production environments, occasional breaking changes can be expected.

> [!NOTE]
> Bugs and issues are handled solely by Engineering teams. On-call support or SLAs are not available.

## License

[Apache 2.0 License](./LICENSE)