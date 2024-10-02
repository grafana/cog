<?php

namespace Grafana\Foundation\VariantPanelcfgFull;

final class VariantConfig
{
    public static function get(): \Grafana\Foundation\Cog\PanelcfgConfig
    {
        return new \Grafana\Foundation\Cog\PanelcfgConfig(
            identifier: 'timeseries',
            optionsFromArray: [\Grafana\Foundation\VariantPanelcfgFull\Options::class, 'fromArray'],
            fieldConfigFromArray: [\Grafana\Foundation\VariantPanelcfgFull\FieldConfig::class, 'fromArray'],
        );
    }
}