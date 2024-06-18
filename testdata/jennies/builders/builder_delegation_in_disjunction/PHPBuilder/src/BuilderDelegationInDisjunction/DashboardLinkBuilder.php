<?php

namespace Grafana\Foundation\BuilderDelegationInDisjunction;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink>
 */
class DashboardLinkBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink();
    }

    /**
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
