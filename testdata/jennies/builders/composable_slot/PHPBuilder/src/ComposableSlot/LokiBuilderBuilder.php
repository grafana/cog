<?php

namespace Grafana\Foundation\ComposableSlot;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\ComposableSlot\Dashboard>
 */
class LokiBuilderBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\ComposableSlot\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\ComposableSlot\Dashboard();
    }

    /**
     * @return \Grafana\Foundation\ComposableSlot\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\Runtime\Dataquery> $target
     */
    public function target(\Grafana\Foundation\Runtime\Builder $target): static
    {
        $targetResource = $target->build();
        $this->internal->target = $targetResource;
    
        return $this;
    }
    /**
     * @param array<\Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\Runtime\Dataquery>> $targets
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
