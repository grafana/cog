<?php

namespace Grafana\Foundation\Sandbox;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Sandbox\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Sandbox\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Sandbox\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Sandbox\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function editable(): static
    {
        $this->internal->editable = true;
    
        return $this;
    }

    public function readonly(): static
    {
        $this->internal->editable = false;
    
        return $this;
    }

    public function autoRefresh(): static
    {
        $this->internal->autoRefresh = true;
    
        return $this;
    }

    public function noAutoRefresh(): static
    {
        $this->internal->autoRefresh = false;
    
        return $this;
    }

}
