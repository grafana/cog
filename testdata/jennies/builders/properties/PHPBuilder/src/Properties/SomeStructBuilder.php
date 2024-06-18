<?php

namespace Grafana\Foundation\Properties;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\Properties\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\Properties\SomeStruct $internal;
    private string $someBuilderProperty;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Properties\SomeStruct();
    }

    /**
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
