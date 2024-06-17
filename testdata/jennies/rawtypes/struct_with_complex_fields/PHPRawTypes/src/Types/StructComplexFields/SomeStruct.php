<?php

namespace Grafana\Foundation\Types\StructComplexFields;

/**
 * This struct does things.
 */
class SomeStruct implements \JsonSerializable {
    public \Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct $fieldRef;

    /**
     * @var string|bool
     */
    public $fieldDisjunctionOfScalars;

    /**
     * @var string|\Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct
     */
    public $fieldMixedDisjunction;

    public ?string $fieldDisjunctionWithNull;

    public \Grafana\Foundation\Types\StructComplexFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public array $fieldArrayOfStrings;

    /**
     * @var array<string, string>
     */
    public array $fieldMapOfStringToString;

    public \Grafana\Foundation\Types\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

    public string $fieldRefToConstant;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "FieldRef" => $this->fieldRef,
            "FieldDisjunctionOfScalars" => $this->fieldDisjunctionOfScalars,
            "FieldMixedDisjunction" => $this->fieldMixedDisjunction,
            "FieldDisjunctionWithNull" => $this->fieldDisjunctionWithNull,
            "Operator" => $this->operator,
            "FieldArrayOfStrings" => $this->fieldArrayOfStrings,
            "FieldMapOfStringToString" => $this->fieldMapOfStringToString,
            "FieldAnonymousStruct" => $this->fieldAnonymousStruct,
            "fieldRefToConstant" => $this->fieldRefToConstant,
        ];
        return $data;
    }
}
