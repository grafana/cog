<?php

namespace Grafana\Foundation\Properties;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Properties\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Properties\SomeStruct $internal;
    private string $someBuilderProperty;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Properties\SomeStruct();
        $this->someBuilderProperty = "";
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Properties\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function id(int $id): static
    {
        $this->internal->id = $id;
    
        return $this;
    }

}
