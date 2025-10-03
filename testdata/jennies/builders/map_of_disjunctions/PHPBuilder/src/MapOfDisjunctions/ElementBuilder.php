<?php

namespace Grafana\Foundation\MapOfDisjunctions;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\Element>
 */
class ElementBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfDisjunctions\Element $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfDisjunctions\Element();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfDisjunctions\Element
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\Panel> $panel
     */
    public function panel(\Grafana\Foundation\Cog\Builder $panel): static
    {
        $panelResource = $panel->build();
        $this->internal->panel = $panelResource;
    
        return $this;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\LibraryPanel> $libraryPanel
     */
    public function libraryPanel(\Grafana\Foundation\Cog\Builder $libraryPanel): static
    {
        $libraryPanelResource = $libraryPanel->build();
        $this->internal->libraryPanel = $libraryPanelResource;
    
        return $this;
    }

}
