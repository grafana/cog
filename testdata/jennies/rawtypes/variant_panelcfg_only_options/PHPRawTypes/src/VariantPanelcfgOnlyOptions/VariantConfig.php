<?php

namespace Grafana\Foundation\VariantPanelcfgOnlyOptions;

final class VariantConfig
{
    public static function get(): \Grafana\Foundation\Cog\PanelcfgConfig
    {
        return new \Grafana\Foundation\Cog\PanelcfgConfig(
            identifier: 'text',
            optionsFromArray: [\Grafana\Foundation\VariantPanelcfgOnlyOptions\Options::class, 'fromArray'],
            fieldConfigFromArray: null,
        );
    }
}