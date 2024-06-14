<?php

namespace Types\StructOptionalFields;

class SomeStruct
{
    public ?\Types\StructOptionalFields\SomeOtherStruct $fieldRef;

    public ?string $fieldString;

    public ?\Types\StructOptionalFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public ?array $fieldArrayOfStrings;

    public ?\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;
}
