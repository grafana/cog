<?php

namespace Grafana\Foundation\MapOfDisjunctions;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\Dashboard>
 */
class DashboardBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfDisjunctions\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfDisjunctions\Dashboard();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfDisjunctions\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param array<string, \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\Element>> $panels
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
