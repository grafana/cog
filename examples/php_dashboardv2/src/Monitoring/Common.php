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
use Grafana\Foundation\Gauge;
use Grafana\Foundation\Logs;
use Grafana\Foundation\Loki;
use Grafana\Foundation\Timeseries;
use Grafana\Foundation\Prometheus;

class Common
{
    public static function defaultTimeseries(): Timeseries\VisualizationBuilder
    {
        return (new Timeseries\VisualizationBuilder())
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

    public static function defaultLogs(): Logs\VisualizationBuilder
    {
        return (new Logs\VisualizationBuilder())
            ->showTime(true)
            ->enableLogDetails(true)
            ->sortOrder(LogsSortOrder::descending())
            ->wrapLogMessage(true);
    }

    public static function defaultGauge(): Gauge\VisualizationBuilder
    {
        return (new Gauge\VisualizationBuilder())
            ->orientation(VizOrientation::auto())
            ->reduceOptions(
                (new ReduceDataOptionsBuilder())
                    ->calcs(['lastNotNull'])
                    ->values(false)
            );
    }

    public static function basicPrometheusQuery(string $query, string $legend): Prometheus\QueryBuilder
    {
        return (new Prometheus\QueryBuilder())
            ->expr($query)
            ->legendFormat($legend);
    }


    public static function tablePrometheusQuery(string $query): Prometheus\QueryBuilder
    {
        return (new Prometheus\QueryBuilder())
            ->expr($query)
            ->instant()
            ->format(Prometheus\PromQueryFormat::table());
    }

    public static function basicLokiQuery(string $query): Loki\QueryBuilder
    {
        return (new Loki\QueryBuilder())
            ->expr($query);
    }
}
