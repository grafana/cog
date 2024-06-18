<?php

namespace Grafana\Foundation\BuilderDelegationInDisjunction;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink>
 */
class ExternalLinkBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink();
    }

    /**
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
