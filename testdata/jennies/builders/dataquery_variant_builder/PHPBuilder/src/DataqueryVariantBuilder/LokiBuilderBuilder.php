<?php

namespace Grafana\Foundation\DataqueryVariantBuilder;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\DataqueryVariantBuilder\Loki>
 */
class LokiBuilderBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\DataqueryVariantBuilder\Loki $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\DataqueryVariantBuilder\Loki();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\DataqueryVariantBuilder\Loki
     */
    public function build()
    {
        return $this->internal;
    }

    public function expr(string $expr): static
    {
        $this->internal->expr = $expr;
    
        return $this;
    }

}
