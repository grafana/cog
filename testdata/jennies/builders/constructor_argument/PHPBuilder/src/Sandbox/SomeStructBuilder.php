<?php

namespace Grafana\Foundation\Sandbox;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Sandbox\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Sandbox\SomeStruct $internal;

    public function __construct(string $title)
    {
    	$this->internal = new \Grafana\Foundation\Sandbox\SomeStruct();
    $this->internal->title = $title;
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Sandbox\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function title(string $title): static
    {
        $this->internal->title = $title;
    
        return $this;
    }

}
