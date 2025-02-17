<?php

namespace Grafana\Foundation\BuilderDelegationInDisjunction;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink>
 */
class DashboardLinkBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink
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
