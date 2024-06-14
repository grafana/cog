<?php

namespace Types\StructComplexFields;

/**
 * This struct does things.
 */
class SomeStruct
{
    public \Types\StructComplexFields\SomeOtherStruct $fieldRef;

    /**
     * @var string|bool
     */
    public $fieldDisjunctionOfScalars;

    /**
     * @var string|\Types\StructComplexFields\SomeOtherStruct
     */
    public $fieldMixedDisjunction;

    public ?string $fieldDisjunctionWithNull;

    public \Types\StructComplexFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public array $fieldArrayOfStrings;

    /**
     * @var array<string, string>
     */
    public array $fieldMapOfStringToString;

    public \Types\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

    public string $fieldRefToConstant;
}
