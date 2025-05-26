<?php

namespace Grafana\Foundation\StructComplexFields;

/**
 * This struct does things.
 */
class SomeStruct implements \JsonSerializable
{
    public \Grafana\Foundation\StructComplexFields\SomeOtherStruct $fieldRef;

    public \Grafana\Foundation\StructComplexFields\StringOrBool $fieldDisjunctionOfScalars;

    public \Grafana\Foundation\StructComplexFields\StringOrSomeOtherStruct $fieldMixedDisjunction;

    public ?string $fieldDisjunctionWithNull;

    public \Grafana\Foundation\StructComplexFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public array $fieldArrayOfStrings;

    /**
     * @var array<string, string>
     */
    public array $fieldMapOfStringToString;

    public \Grafana\Foundation\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

    public string $fieldRefToConstant;

    /**
     * @param \Grafana\Foundation\StructComplexFields\SomeOtherStruct|null $fieldRef
     * @param \Grafana\Foundation\StructComplexFields\StringOrBool|null $fieldDisjunctionOfScalars
     * @param \Grafana\Foundation\StructComplexFields\StringOrSomeOtherStruct|null $fieldMixedDisjunction
     * @param string|null $fieldDisjunctionWithNull
     * @param \Grafana\Foundation\StructComplexFields\SomeStructOperator|null $operator
     * @param array<string>|null $fieldArrayOfStrings
     * @param array<string, string>|null $fieldMapOfStringToString
     * @param \Grafana\Foundation\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct|null $fieldAnonymousStruct
     * @param \Grafana\Foundation\StructComplexFields\ConnectionPath|null $fieldRefToConstant
     */
    public function __construct(?\Grafana\Foundation\StructComplexFields\SomeOtherStruct $fieldRef = null, ?\Grafana\Foundation\StructComplexFields\StringOrBool $fieldDisjunctionOfScalars = null, ?\Grafana\Foundation\StructComplexFields\StringOrSomeOtherStruct $fieldMixedDisjunction = null, ?string $fieldDisjunctionWithNull = null, ?\Grafana\Foundation\StructComplexFields\SomeStructOperator $operator = null, ?array $fieldArrayOfStrings = null, ?array $fieldMapOfStringToString = null, ?\Grafana\Foundation\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct = null, ?\Grafana\Foundation\StructComplexFields\ConnectionPath $fieldRefToConstant = null)
    {
        $this->fieldRef = $fieldRef ?: new \Grafana\Foundation\StructComplexFields\SomeOtherStruct();
        $this->fieldDisjunctionOfScalars = $fieldDisjunctionOfScalars ?: new \Grafana\Foundation\StructComplexFields\StringOrBool();
        $this->fieldMixedDisjunction = $fieldMixedDisjunction ?: new \Grafana\Foundation\StructComplexFields\StringOrSomeOtherStruct();
        $this->fieldDisjunctionWithNull = $fieldDisjunctionWithNull;
        $this->operator = $operator ?: \Grafana\Foundation\StructComplexFields\SomeStructOperator::GreaterThan();
        $this->fieldArrayOfStrings = $fieldArrayOfStrings ?: [];
        $this->fieldMapOfStringToString = $fieldMapOfStringToString ?: [];
        $this->fieldAnonymousStruct = $fieldAnonymousStruct ?: new \Grafana\Foundation\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct();
        $this->fieldRefToConstant = $fieldRefToConstant ?: \Grafana\Foundation\StructComplexFields\ConnectionPath;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldRef?: mixed, FieldDisjunctionOfScalars?: mixed, FieldMixedDisjunction?: mixed, FieldDisjunctionWithNull?: string, Operator?: string, FieldArrayOfStrings?: array<string>, FieldMapOfStringToString?: array<string, string>, FieldAnonymousStruct?: mixed, fieldRefToConstant?: string} $inputData */
        $data = $inputData;
        return new self(
            fieldRef: isset($data["FieldRef"]) ? (function($input) {
    	/** @var array{FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructComplexFields\SomeOtherStruct::fromArray($val);
    })($data["FieldRef"]) : null,
            fieldDisjunctionOfScalars: isset($data["FieldDisjunctionOfScalars"]) ? (function($input) {
    	/** @var array{String?: string, Bool?: bool} */
    $val = $input;
    	return \Grafana\Foundation\StructComplexFields\StringOrBool::fromArray($val);
    })($data["FieldDisjunctionOfScalars"]) : null,
            fieldMixedDisjunction: isset($data["FieldMixedDisjunction"]) ? (function($input) {
    	/** @var array{String?: string, SomeOtherStruct?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructComplexFields\StringOrSomeOtherStruct::fromArray($val);
    })($data["FieldMixedDisjunction"]) : null,
            fieldDisjunctionWithNull: $data["FieldDisjunctionWithNull"] ?? null,
            operator: isset($data["Operator"]) ? (function($input) { return \Grafana\Foundation\StructComplexFields\SomeStructOperator::fromValue($input); })($data["Operator"]) : null,
            fieldArrayOfStrings: $data["FieldArrayOfStrings"] ?? null,
            fieldMapOfStringToString: $data["FieldMapOfStringToString"] ?? null,
            fieldAnonymousStruct: isset($data["FieldAnonymousStruct"]) ? (function($input) {
    	/** @var array{FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructComplexFields\StructComplexFieldsSomeStructFieldAnonymousStruct::fromArray($val);
    })($data["FieldAnonymousStruct"]) : null,
            fieldRefToConstant: isset($data["fieldRefToConstant"]) ? /* ref to a non-struct, non-enum, this should have been inlined */ (function(array $input) { return $input; })($data["fieldRefToConstant"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->FieldRef = $this->fieldRef;
        $data->FieldDisjunctionOfScalars = $this->fieldDisjunctionOfScalars;
        $data->FieldMixedDisjunction = $this->fieldMixedDisjunction;
        $data->Operator = $this->operator;
        $data->FieldArrayOfStrings = $this->fieldArrayOfStrings;
        $data->FieldMapOfStringToString = $this->fieldMapOfStringToString;
        $data->FieldAnonymousStruct = $this->fieldAnonymousStruct;
        $data->fieldRefToConstant = $this->fieldRefToConstant;
        if (isset($this->fieldDisjunctionWithNull)) {
            $data->FieldDisjunctionWithNull = $this->fieldDisjunctionWithNull;
        }
        return $data;
    }
}
