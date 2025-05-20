<?php

namespace Grafana\Foundation\StructOptionalFields;

class SomeStruct implements \JsonSerializable
{
    public ?\Grafana\Foundation\StructOptionalFields\SomeOtherStruct $fieldRef;

    public ?string $fieldString;

    public ?\Grafana\Foundation\StructOptionalFields\SomeStructOperator $operator;

    /**
     * @var array<string>|null
     */
    public ?array $fieldArrayOfStrings;

    public ?\Grafana\Foundation\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

    /**
     * @param \Grafana\Foundation\StructOptionalFields\SomeOtherStruct|null $fieldRef
     * @param string|null $fieldString
     * @param \Grafana\Foundation\StructOptionalFields\SomeStructOperator|null $operator
     * @param array<string>|null $fieldArrayOfStrings
     * @param \Grafana\Foundation\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct|null $fieldAnonymousStruct
     */
    public function __construct(?\Grafana\Foundation\StructOptionalFields\SomeOtherStruct $fieldRef = null, ?string $fieldString = null, ?\Grafana\Foundation\StructOptionalFields\SomeStructOperator $operator = null, ?array $fieldArrayOfStrings = null, ?\Grafana\Foundation\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct = null)
    {
        $this->fieldRef = $fieldRef;
        $this->fieldString = $fieldString;
        $this->operator = $operator;
        $this->fieldArrayOfStrings = $fieldArrayOfStrings;
        $this->fieldAnonymousStruct = $fieldAnonymousStruct;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldRef?: mixed, FieldString?: string, Operator?: string, FieldArrayOfStrings?: array<string>, FieldAnonymousStruct?: mixed} $inputData */
        $data = $inputData;
        return new self(
            fieldRef: isset($data["FieldRef"]) ? (function($input) {
    	/** @var array{FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructOptionalFields\SomeOtherStruct::fromArray($val);
    })($data["FieldRef"]) : null,
            fieldString: $data["FieldString"] ?? null,
            operator: isset($data["Operator"]) ? (function($input) { return \Grafana\Foundation\StructOptionalFields\SomeStructOperator::fromValue($input); })($data["Operator"]) : null,
            fieldArrayOfStrings: $data["FieldArrayOfStrings"] ?? null,
            fieldAnonymousStruct: isset($data["FieldAnonymousStruct"]) ? (function($input) {
    	/** @var array{FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct::fromArray($val);
    })($data["FieldAnonymousStruct"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->fieldRef)) {
            $data->FieldRef = $this->fieldRef;
        }
        if (isset($this->fieldString)) {
            $data->FieldString = $this->fieldString;
        }
        if (isset($this->operator)) {
            $data->Operator = $this->operator;
        }
        if (isset($this->fieldArrayOfStrings)) {
            $data->FieldArrayOfStrings = $this->fieldArrayOfStrings;
        }
        if (isset($this->fieldAnonymousStruct)) {
            $data->FieldAnonymousStruct = $this->fieldAnonymousStruct;
        }
        return $data;
    }
}
