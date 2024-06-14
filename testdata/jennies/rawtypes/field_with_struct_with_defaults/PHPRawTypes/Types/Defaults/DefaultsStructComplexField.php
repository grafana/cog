<?php

namespace Types\Defaults;

class DefaultsStructComplexField
{
    public string $uid;

    public \Types\Defaults\DefaultsStructComplexFieldNested $nested;

    /**
     * @var array<string>
     */
    public array $array;
}
