<?php

namespace Grafana\Foundation\Sandbox;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Sandbox\Dashboard>
 */
class DashboardBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Sandbox\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Sandbox\Dashboard();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Sandbox\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    public function withVariable(string $name, string $value): static
    {
        $this->internal->variables[] = new \Grafana\Foundation\Sandbox\Variable(
            name: $name,
            value: $value,
        );
    
        return $this;
    }

}
