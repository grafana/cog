<?php

namespace Grafana\Foundation\Sandbox;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\Sandbox\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\Sandbox\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Sandbox\SomeStruct();
    }

    /**
     * @return \Grafana\Foundation\Sandbox\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function tags(string $tags): static
    {
        $this->internal->tags[] = $tags;
    
        return $this;
    }

}
