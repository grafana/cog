<?php

namespace Grafana\Foundation\Sandbox;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Sandbox\SomeStructWithDefaultEnum>
 */
class SomeStructWithDefaultEnumBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Sandbox\SomeStructWithDefaultEnum $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Sandbox\SomeStructWithDefaultEnum();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Sandbox\SomeStructWithDefaultEnum
     */
    public function build()
    {
        return $this->internal;
    }

    public function data(\Grafana\Foundation\Sandbox\StringEnumWithDefault $key, string $value): static
    {
        $this->internal->data[$key] = $value;
    
        return $this;
    }

}
