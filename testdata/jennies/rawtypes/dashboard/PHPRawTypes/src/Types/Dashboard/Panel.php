<?php

namespace Grafana\Foundation\Types\Dashboard;

class Panel
{
    public string $title;

    public string $type;

    public ?\Grafana\Foundation\Types\Dashboard\DataSourceRef $datasource;

    /**
     * @var mixed
     */
    public $options;

    /**
     * @var array<\Grafana\Foundation\Runtime\Variants\Dataquery>
     */
    public ?array $targets;

    public ?\Grafana\Foundation\Types\Dashboard\FieldConfigSource $fieldConfig;
}
