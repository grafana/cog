<?php

namespace Grafana\Foundation\Builderpkg;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\Withdashes\SomeStruct>
 */
class SomeNiceBuilderBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\Withdashes\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Withdashes\SomeStruct();
    }

    /**
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
