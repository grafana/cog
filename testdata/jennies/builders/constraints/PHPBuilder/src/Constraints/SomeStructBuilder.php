<?php

namespace Grafana\Foundation\Constraints;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Constraints\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Constraints\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Constraints\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Constraints\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    public function id(int $id): static
    {
        if (!($id >= 5)) {
            throw new \ValueError('$id must be >= 5');
        }
        if (!($id < 10)) {
            throw new \ValueError('$id must be < 10');
        }
        $this->internal->id = $id;
    
        return $this;
    }

    public function title(string $title): static
    {
        if (!(strlen($title) >= 1)) {
            throw new \ValueError('strlen($title) must be >= 1');
        }
        $this->internal->title = $title;
    
        return $this;
    }

}
