<?php

namespace Grafana\Foundation\AnonymousStruct;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\AnonymousStruct\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\AnonymousStruct\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\AnonymousStruct\SomeStruct();
    }

    /**
     * @return \Grafana\Foundation\AnonymousStruct\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function time(unknown $time): static
    {
        $this->internal->time = $time;
    
        return $this;
    }

}
