<?php

namespace Grafana\Foundation\InitializationSafeguards;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\InitializationSafeguards\SomePanel>
 */
class SomePanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\InitializationSafeguards\SomePanel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\InitializationSafeguards\SomePanel();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\InitializationSafeguards\SomePanel
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

    public function showLegend(bool $show): static
    {    
        if ($this->internal->options === null) {
            $this->internal->options = new \Grafana\Foundation\InitializationSafeguards\Options();
        }
        assert($this->internal->options instanceof \Grafana\Foundation\InitializationSafeguards\Options);
        $this->internal->options->legend->show = $show;
    
        return $this;
    }

}
