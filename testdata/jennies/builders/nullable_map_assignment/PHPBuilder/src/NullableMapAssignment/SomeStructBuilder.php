<?php

namespace Grafana\Foundation\NullableMapAssignment;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\NullableMapAssignment\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\NullableMapAssignment\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\NullableMapAssignment\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\NullableMapAssignment\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param array<string, string> $config
     */
    public function config(array $config): static
    {
        $this->internal->config = $config;
    
        return $this;
    }

}
