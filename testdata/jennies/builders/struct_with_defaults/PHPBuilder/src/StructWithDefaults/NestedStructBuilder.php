<?php

namespace Grafana\Foundation\StructWithDefaults;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\StructWithDefaults\NestedStruct>
 */
class NestedStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\StructWithDefaults\NestedStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\StructWithDefaults\NestedStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\StructWithDefaults\NestedStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function stringVal(string $stringVal): static
    {
        $this->internal->stringVal = $stringVal;
    
        return $this;
    }

    public function intVal(int $intVal): static
    {
        $this->internal->intVal = $intVal;
    
        return $this;
    }

}
