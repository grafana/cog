<?php

namespace Grafana\Foundation\Panelbuilder;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Panelbuilder\Panel>
 */
class PanelBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Panelbuilder\Panel $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Panelbuilder\Panel();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Panelbuilder\Panel
     */
    public function build()
    {
        return $this->internal;
    }

    public function onlyFromThisDashboard(bool $onlyFromThisDashboard): static
    {
        $this->internal->onlyFromThisDashboard = $onlyFromThisDashboard;
    
        return $this;
    }

    public function onlyInTimeRange(bool $onlyInTimeRange): static
    {
        $this->internal->onlyInTimeRange = $onlyInTimeRange;
    
        return $this;
    }

    /**
     * @param array<string> $tags
     */
    public function tags(array $tags): static
    {
        $this->internal->tags = $tags;
    
        return $this;
    }

    public function limit(int $limit): static
    {
        $this->internal->limit = $limit;
    
        return $this;
    }

    public function showUser(bool $showUser): static
    {
        $this->internal->showUser = $showUser;
    
        return $this;
    }

    public function showTime(bool $showTime): static
    {
        $this->internal->showTime = $showTime;
    
        return $this;
    }

    public function showTags(bool $showTags): static
    {
        $this->internal->showTags = $showTags;
    
        return $this;
    }

    public function navigateToPanel(bool $navigateToPanel): static
    {
        $this->internal->navigateToPanel = $navigateToPanel;
    
        return $this;
    }

    public function navigateBefore(string $navigateBefore): static
    {
        $this->internal->navigateBefore = $navigateBefore;
    
        return $this;
    }

    public function navigateAfter(string $navigateAfter): static
    {
        $this->internal->navigateAfter = $navigateAfter;
    
        return $this;
    }

}
