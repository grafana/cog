<?php

namespace Grafana\Foundation\MapOfDisjunctions;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\Panel>
 */
class PanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfDisjunctions\Panel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfDisjunctions\Panel();
    $this->internal->kind = "Panel";
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfDisjunctions\Panel
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
