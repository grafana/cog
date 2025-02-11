<?php

namespace Grafana\Foundation\BuilderDelegation;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegation\DashboardLink>
 */
class DashboardLinkBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BuilderDelegation\DashboardLink $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegation\DashboardLink();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BuilderDelegation\DashboardLink
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

    public function url(string $url): static
    {
        $this->internal->url = $url;
    
        return $this;
    }

}
