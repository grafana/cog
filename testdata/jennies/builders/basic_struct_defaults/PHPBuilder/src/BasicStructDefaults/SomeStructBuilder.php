<?php

namespace Grafana\Foundation\BasicStructDefaults;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BasicStructDefaults\SomeStruct>
 */
class SomeStructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BasicStructDefaults\SomeStruct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BasicStructDefaults\SomeStruct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BasicStructDefaults\SomeStruct
     */
    public function build()
    {
        return $this->internal;
    }

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

    public function liveNow(bool $liveNow): static
    {
        $this->internal->liveNow = $liveNow;
    
        return $this;
    }

}
