<?php

namespace Grafana\Foundation\Builderpkg;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Withdashes\SomeStruct>
 */
class SomeNiceBuilderBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Withdashes\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Withdashes\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Withdashes\SomeStruct
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
