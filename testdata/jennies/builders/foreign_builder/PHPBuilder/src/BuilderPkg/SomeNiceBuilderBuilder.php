<?php

namespace Grafana\Foundation\BuilderPkg;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\SomePkg\SomeStruct>
 */
class SomeNiceBuilderBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\SomePkg\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\SomePkg\SomeStruct();
    }

    /**
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
