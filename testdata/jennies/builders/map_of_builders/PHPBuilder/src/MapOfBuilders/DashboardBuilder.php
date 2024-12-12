<?php

namespace Grafana\Foundation\MapOfBuilders;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfBuilders\Dashboard>
 */
class DashboardBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfBuilders\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfBuilders\Dashboard();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfBuilders\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param array<string, \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfBuilders\Panel>> $panels
     */
    public function panels(array $panels): static
    {
            $panelsResources = [];
            foreach ($panels as $key1 => $val1) {
                    $panelsResources[$key1] = $val1->build();
            }
        $this->internal->panels = $panelsResources;
    
        return $this;
    }

}
