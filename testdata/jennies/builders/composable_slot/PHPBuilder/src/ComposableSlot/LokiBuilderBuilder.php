<?php

namespace Grafana\Foundation\ComposableSlot;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\ComposableSlot\Dashboard>
 */
class LokiBuilderBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\ComposableSlot\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\ComposableSlot\Dashboard();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\ComposableSlot\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Cog\Dataquery> $target
     */
    public function target(\Grafana\Foundation\Cog\Builder $target): static
    {
        $targetResource = $target->build();
        $this->internal->target = $targetResource;
    
        return $this;
    }

    /**
     * @param array<\Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Cog\Dataquery>> $targets
     */
    public function targets(array $targets): static
    {
            $targetsResources = [];
            foreach ($targets as $r1) {
                    $targetsResources[] = $r1->build();
            }
        $this->internal->targets = $targetsResources;
    
        return $this;
    }

}
