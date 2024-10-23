<?php

namespace Grafana\Foundation\BuilderDelegationInDisjunction;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink>
 */
class ExternalLinkBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink
     */
    public function build()
    {
        return $this->internal;
    }

    public function url(string $url): static
    {
        $this->internal->url = $url;
    
        return $this;
    }

}
