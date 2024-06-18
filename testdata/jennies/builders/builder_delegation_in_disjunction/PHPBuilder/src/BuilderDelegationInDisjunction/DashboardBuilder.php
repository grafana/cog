<?php

namespace Grafana\Foundation\BuilderDelegationInDisjunction;

/**
 * @implements \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\Dashboard>
 */
class DashboardBuilder implements \Grafana\Foundation\Runtime\Builder
{
    protected \Grafana\Foundation\BuilderDelegationInDisjunction\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegationInDisjunction\Dashboard();
    }

    /**
     * @return \Grafana\Foundation\BuilderDelegationInDisjunction\Dashboard
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * will be expanded to cog.Builder<DashboardLink> | string
     * @param \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink>|string $singleLinkOrString
     */
    public function singleLinkOrString( $singleLinkOrString): static
    {
        /** @var \Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink|string $singleLinkOrStringResource */
        $singleLinkOrStringResource = $singleLinkOrString instanceof \Grafana\Foundation\Runtime\Builder ? $singleLinkOrString->build() : $singleLinkOrString;
        $this->internal->singleLinkOrString = $singleLinkOrStringResource;
    
        return $this;
    }
    /**
     * will be expanded to [](cog.Builder<DashboardLink> | string)
     * @param array<\Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink>|string> $linksOrStrings
     */
    public function linksOrStrings(array $linksOrStrings): static
    {
            $linksOrStringsResources = [];
            foreach ($linksOrStrings as $r1) {
                    $linksOrStringsResources[] = $r1 instanceof \Grafana\Foundation\Runtime\Builder ? $r1->build() : $r1;
            }
        $this->internal->linksOrStrings = $linksOrStringsResources;
    
        return $this;
    }
    /**
     * @param \Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\DashboardLink>|\Grafana\Foundation\Runtime\Builder<\Grafana\Foundation\BuilderDelegationInDisjunction\ExternalLink> $disjunctionOfBuilders
     */
    public function disjunctionOfBuilders( $disjunctionOfBuilders): static
    {
        $disjunctionOfBuildersResource = $disjunctionOfBuilders->build();
        $this->internal->disjunctionOfBuilders = $disjunctionOfBuildersResource;
    
        return $this;
    }

}
