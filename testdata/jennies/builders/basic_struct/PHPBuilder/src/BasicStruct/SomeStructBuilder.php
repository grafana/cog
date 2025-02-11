<?php

namespace Grafana\Foundation\BasicStruct;

/**
 * SomeStruct, to hold data.
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BasicStruct\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BasicStruct\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BasicStruct\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BasicStruct\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * id identifies something. Weird, right?
     */
    public function id(int $id): static
    {
        $this->internal->id = $id;
    
        return $this;
    }

    public function uid(string $uid): static
    {
        $this->internal->uid = $uid;
    
        return $this;
    }

    /**
     * @param array<string> $tags
     */
    public function tags(array $tags): static
    {
        $this->internal->tags = $tags;
    
        return $this;
    }

    /**
     * This thing could be live.
     * Or maybe not.
     */
    public function liveNow(bool $liveNow): static
    {
        $this->internal->liveNow = $liveNow;
    
        return $this;
    }

}
