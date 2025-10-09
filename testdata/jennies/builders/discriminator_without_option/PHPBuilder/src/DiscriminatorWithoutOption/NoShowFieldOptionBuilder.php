<?php

namespace Grafana\Foundation\DiscriminatorWithoutOption;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOption>
 */
class NoShowFieldOptionBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOption $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOption();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOption
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
