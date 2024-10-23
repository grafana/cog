<?php

namespace Grafana\Foundation\BuilderPkg;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\SomePkg\SomeStruct>
 */
class SomeNiceBuilderBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\SomePkg\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\SomePkg\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\SomePkg\SomeStruct
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
