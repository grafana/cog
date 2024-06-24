<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\GraphDrawStyle;
use Grafana\Foundation\Common\LegendDisplayMode;
use Grafana\Foundation\Common\LegendPlacement;
use Grafana\Foundation\Common\LogsSortOrder;
use Grafana\Foundation\Common\ReduceDataOptionsBuilder;
use Grafana\Foundation\Common\VisibilityMode;
use Grafana\Foundation\Common\VizLegendOptionsBuilder;
use Grafana\Foundation\Common\VizOrientation;
use Grafana\Foundation\Dashboard\DataSourceRef;
use Grafana\Foundation\Gauge;
use Grafana\Foundation\Logs;
use Grafana\Foundation\Loki;
use Grafana\Foundation\Prometheus\PromQueryFormat;
use Grafana\Foundation\Timeseries;
use Grafana\Foundation\Prometheus;

class Common
{
    public static function defaultTimeseries(): Timeseries\PanelBuilder
    {
        return (new Timeseries\PanelBuilder())
            ->lineWidth(1)
            ->fillOpacity(0)
            ->drawStyle(GraphDrawStyle::line())
            ->showPoints(VisibilityMode::never())
            ->legend(
                (new VizLegendOptionsBuilder())
                    ->showLegend(true)
                    ->placement(LegendPlacement::bottom())
                    ->displayMode(LegendDisplayMode::list())
            );
    }

    public static function defaultLogs(): Logs\PanelBuilder
    {
        return (new Logs\PanelBuilder())
            ->span(24)
            ->datasource(new DataSourceRef(type: 'loki', uid: 'grafanacloud-logs'))
            ->showTime(true)
            ->enableLogDetails(true)
            ->sortOrder(LogsSortOrder::descending())
            ->wrapLogMessage(true);
    }

    public static function defaultGauge(): Gauge\PanelBuilder
    {
        return (new Gauge\PanelBuilder())
            ->orientation(VizOrientation::auto())
            ->reduceOptions(
                (new ReduceDataOptionsBuilder())
                    ->calcs(['lastNotNull'])
                    ->values(false)
            );
    }

    public static function basicPrometheusQuery(string $query, string $legend): Prometheus\DataqueryBuilder
    {
        return (new Prometheus\DataqueryBuilder())
            ->expr($query)
            ->legendFormat($legend);
    }


    public static function tablePrometheusQuery(string $query, string $ref): Prometheus\DataqueryBuilder
    {
        return (new Prometheus\DataqueryBuilder())
            ->expr($query)
            ->instant()
            ->format(PromQueryFormat::table())
            ->refId($ref);
    }

    public static function basicLokiQuery(string $query): Loki\DataqueryBuilder
    {
        return (new Loki\DataqueryBuilder())
            ->expr($query);
    }
}