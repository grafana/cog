<?php

namespace Grafana\Foundation\Types\Dashboard;

class Dashboard
{
    public string $title;

    /**
     * @var array<\Grafana\Foundation\Types\Dashboard\Panel>
     */
    public ?array $panels;
}
