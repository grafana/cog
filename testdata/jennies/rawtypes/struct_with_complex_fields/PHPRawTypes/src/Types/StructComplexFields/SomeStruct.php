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
     * @param \Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct|null $fieldRef
     * @param string|bool|null $fieldDisjunctionOfScalars
     * @param string|\Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct|null $fieldMixedDisjunction
     * @param string|null $fieldDisjunctionWithNull
     * @param \Grafana\Foundation\Types\StructComplexFields\SomeStructOperator|null $operator
     * @param array<string>|null $fieldArrayOfStrings
     * @param array<string, string>|null $fieldMapOfStringToString
     * @param \Grafana\Foundation\Types\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct|null $fieldAnonymousStruct
     * @param \Grafana\Foundation\Types\StructComplexFields\ConnectionPath|null $fieldRefToConstant
     */
    public function __construct(?\Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct $fieldRef = null,  $fieldDisjunctionOfScalars = null,  $fieldMixedDisjunction = null, ?string $fieldDisjunctionWithNull = null, ?\Grafana\Foundation\Types\StructComplexFields\SomeStructOperator $operator = null, ?array $fieldArrayOfStrings = null, ?array $fieldMapOfStringToString = null, ?\Grafana\Foundation\Types\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct = null, ?\Grafana\Foundation\Types\StructComplexFields\ConnectionPath $fieldRefToConstant = null)
    {
        $this->fieldRef = $fieldRef ?: new \Grafana\Foundation\Types\StructComplexFields\SomeOtherStruct();
        $this->fieldDisjunctionOfScalars = $fieldDisjunctionOfScalars ?: "";
        $this->fieldMixedDisjunction = $fieldMixedDisjunction ?: "";
        $this->fieldDisjunctionWithNull = $fieldDisjunctionWithNull;
        $this->operator = $operator ?: \Grafana\Foundation\Types\StructComplexFields\SomeStructOperator::GreaterThan();
        $this->fieldArrayOfStrings = $fieldArrayOfStrings ?: [];
        $this->fieldMapOfStringToString = $fieldMapOfStringToString ?: [];
        $this->fieldAnonymousStruct = $fieldAnonymousStruct ?: new \Grafana\Foundation\Types\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct();
        $this->fieldRefToConstant = $fieldRefToConstant ?: \Grafana\Foundation\Types\StructComplexFields\ConnectionPath;
    }

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
