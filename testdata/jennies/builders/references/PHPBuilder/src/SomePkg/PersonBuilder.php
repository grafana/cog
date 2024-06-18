<?php

namespace Grafana\Foundation\SomePkg;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\SomePkg\Person>
 */
class PersonBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\SomePkg\Person $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\SomePkg\Person();
    }

    /**
     * @return \Grafana\Foundation\SomePkg\Person
     */
    public function build()
    {
        return $this->internal;
    }

    public function name(\Grafana\Foundation\OtherPkg\Name $name): static
    {
        $this->internal->name = $name;
    
        return $this;
    }

}
