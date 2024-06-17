<?php

namespace Grafana\Foundation\Types\Defaults;

class DefaultsStructComplexField
{
    public string $uid;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested $nested;

    /**
     * @var array<string>
     */
    public array $array;
}
