<?php

namespace Types\Dashboard;

class Panel
{
    public string $title;

    public string $type;

    public ?\Types\Dashboard\DataSourceRef $datasource;

    /**
     * @var mixed
     */
    public $options;

    /**
     * @var array<\Runtime\Variants\Dataquery>
     */
    public ?array $targets;

    public ?\Types\Dashboard\FieldConfigSource $fieldConfig;
}
