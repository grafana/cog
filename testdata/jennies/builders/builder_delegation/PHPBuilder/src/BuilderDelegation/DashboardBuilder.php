<?php

namespace Grafana\Foundation\BuilderDelegation;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegation\Dashboard>
 */
class DashboardBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\BuilderDelegation\Dashboard $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\BuilderDelegation\Dashboard();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\BuilderDelegation\Dashboard
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

    public function title(string $title): static
    {
        $this->internal->title = $title;
    
        return $this;
    }

    /**
     * will be expanded to []cog.Builder<DashboardLink>
     * @param array<\Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegation\DashboardLink>> $links
     */
    public function links(array $links): static
    {
            $linksResources = [];
            foreach ($links as $r1) {
                    $linksResources[] = $r1->build();
            }
        $this->internal->links = $linksResources;
    
        return $this;
    }

    /**
     * will be expanded to [][]cog.Builder<DashboardLink>
     * @param array<array<\Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegation\DashboardLink>>> $linksOfLinks
     */
    public function linksOfLinks(array $linksOfLinks): static
    {
            $linksOfLinksResources = [];
            foreach ($linksOfLinks as $r1) {
                    $linksOfLinksDepth1 = [];
            foreach ($r1 as $r2) {
                    $linksOfLinksDepth1[] = $r2->build();
            }
    
                    $linksOfLinksResources[] = $linksOfLinksDepth1;
            }
        $this->internal->linksOfLinks = $linksOfLinksResources;
    
        return $this;
    }

    /**
     * will be expanded to cog.Builder<DashboardLink>
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\BuilderDelegation\DashboardLink> $singleLink
     */
    public function singleLink(\Grafana\Foundation\Cog\Builder $singleLink): static
    {
        $singleLinkResource = $singleLink->build();
        $this->internal->singleLink = $singleLinkResource;
    
        return $this;
    }

}
