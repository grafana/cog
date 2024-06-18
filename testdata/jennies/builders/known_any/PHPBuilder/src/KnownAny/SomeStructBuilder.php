<?php

namespace Grafana\Foundation\KnownAny;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\KnownAny\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\KnownAny\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\KnownAny\SomeStruct();
    }

    /**
     * @return \Grafana\Foundation\KnownAny\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function title(string $title): static
    {    
        if ($this->internal->config === null) {
            $this->internal->config = new \Grafana\Foundation\KnownAny\Config();
        }
        assert($this->internal->config instanceof \Grafana\Foundation\KnownAny\Config);
        $this->internal->config->title = $title;
    
        return $this;
    }

}
