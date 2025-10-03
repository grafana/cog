<?php

namespace Grafana\Foundation\MapOfDisjunctions;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfDisjunctions\LibraryPanel>
 */
class LibraryPanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfDisjunctions\LibraryPanel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfDisjunctions\LibraryPanel();
    $this->internal->kind = "Library";
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfDisjunctions\LibraryPanel
     */
    public function build()
    {
        return $this->internal;
    }

    public function text(string $text): static
    {
        $this->internal->text = $text;
    
        return $this;
    }

}
