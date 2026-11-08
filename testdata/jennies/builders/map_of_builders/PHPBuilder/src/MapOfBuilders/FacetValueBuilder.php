<?php

namespace Grafana\Foundation\MapOfBuilders;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\MapOfBuilders\FacetValue>
 */
class FacetValueBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\MapOfBuilders\FacetValue $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\MapOfBuilders\FacetValue();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\MapOfBuilders\FacetValue
     */
    public function build()
    {
        return $this->internal;
    }

    public function value(string $value): static
    {
        $this->internal->value = $value;
    
        return $this;
    }

    public function count(int $count): static
    {
        $this->internal->count = $count;
    
        return $this;
    }

}
