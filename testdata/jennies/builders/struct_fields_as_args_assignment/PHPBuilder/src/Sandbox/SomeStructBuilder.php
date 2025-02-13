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

    public function time(string $from, string $to): static
    {    
        if ($this->internal->time === null) {
            $this->internal->time = "unknown";
        }
        
        $this->internal->time->from = $from;
        $this->internal->time->to = $to;
    
        return $this;
    }

}
