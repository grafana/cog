<?php

namespace Grafana\Foundation\ConstructorInitializations;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\ConstructorInitializations\SomePanel>
 */
class SomePanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\ConstructorInitializations\SomePanel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\ConstructorInitializations\SomePanel();
    $this->internal->type = "panel_type";
    $this->internal->cursor = \Grafana\Foundation\ConstructorInitializations\CursorMode::tooltip();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\ConstructorInitializations\SomePanel
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
