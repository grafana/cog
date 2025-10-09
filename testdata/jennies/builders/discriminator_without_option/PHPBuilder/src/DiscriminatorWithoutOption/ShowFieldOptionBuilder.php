<?php

namespace Grafana\Foundation\DiscriminatorWithoutOption;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOption>
 */
class ShowFieldOptionBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOption $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOption();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOption
     */
    public function build()
    {
        return $this->internal;
    }

    public function field(\Grafana\Foundation\DiscriminatorWithoutOption\AnEnum $field): static
    {
        $this->internal->field = $field;
    
        return $this;
    }

    public function text(string $text): static
    {
        $this->internal->text = $text;
    
        return $this;
    }

}
