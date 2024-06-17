<?php

namespace Grafana\Foundation\Types\StructOptionalFields;

class SomeStruct
{
    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeOtherStruct $fieldRef;

    public ?string $fieldString;

    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public ?array $fieldArrayOfStrings;

    public ?\Grafana\Foundation\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;
}
