<?php

namespace Grafana\Foundation\MapOfBuilders;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfBuilders\Panel>
 */
class PanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfBuilders\Panel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfBuilders\Panel();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfBuilders\Panel
     */
    public function build()
    {
        return $this->internal;
    }

    public function title(string $title): static
    {
        $this->internal->title = $title;
    
        return $this;
    }

}
