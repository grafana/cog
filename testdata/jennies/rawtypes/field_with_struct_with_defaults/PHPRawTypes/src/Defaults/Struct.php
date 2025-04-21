<?php

namespace Grafana\Foundation\Defaults;

class Struct implements \JsonSerializable
{
    public \Grafana\Foundation\Defaults\NestedStruct $allFields;

    public \Grafana\Foundation\Defaults\NestedStruct $partialFields;

    public \Grafana\Foundation\Defaults\NestedStruct $emptyFields;

    public \Grafana\Foundation\Defaults\DefaultsStructComplexField $complexField;

    public \Grafana\Foundation\Defaults\DefaultsStructPartialComplexField $partialComplexField;

    /**
     * @param \Grafana\Foundation\Defaults\NestedStruct|null $allFields
     * @param \Grafana\Foundation\Defaults\NestedStruct|null $partialFields
     * @param \Grafana\Foundation\Defaults\NestedStruct|null $emptyFields
     * @param \Grafana\Foundation\Defaults\DefaultsStructComplexField|null $complexField
     * @param \Grafana\Foundation\Defaults\DefaultsStructPartialComplexField|null $partialComplexField
     */
    public function __construct(?\Grafana\Foundation\Defaults\NestedStruct $allFields = null, ?\Grafana\Foundation\Defaults\NestedStruct $partialFields = null, ?\Grafana\Foundation\Defaults\NestedStruct $emptyFields = null, ?\Grafana\Foundation\Defaults\DefaultsStructComplexField $complexField = null, ?\Grafana\Foundation\Defaults\DefaultsStructPartialComplexField $partialComplexField = null)
    {
        $this->allFields = $allFields ?: new \Grafana\Foundation\Defaults\NestedStruct(intVal: 3, stringVal: "hello");
        $this->partialFields = $partialFields ?: new \Grafana\Foundation\Defaults\NestedStruct(intVal: 3);
        $this->emptyFields = $emptyFields ?: new \Grafana\Foundation\Defaults\NestedStruct();
        $this->complexField = $complexField ?: new \Grafana\Foundation\Defaults\DefaultsStructComplexField(array: ["hello"], nested: new \Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested(nestedVal: "nested"), uid: "myUID");
        $this->partialComplexField = $partialComplexField ?: new \Grafana\Foundation\Defaults\DefaultsStructPartialComplexField();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{allFields?: mixed, partialFields?: mixed, emptyFields?: mixed, complexField?: mixed, partialComplexField?: mixed} $inputData */
        $data = $inputData;
        return new self(
            allFields: isset($data["allFields"]) ? (function($input) {
    	/** @var array{stringVal?: string, intVal?: int} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\NestedStruct::fromArray($val);
    })($data["allFields"]) : null,
            partialFields: isset($data["partialFields"]) ? (function($input) {
    	/** @var array{stringVal?: string, intVal?: int} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\NestedStruct::fromArray($val);
    })($data["partialFields"]) : null,
            emptyFields: isset($data["emptyFields"]) ? (function($input) {
    	/** @var array{stringVal?: string, intVal?: int} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\NestedStruct::fromArray($val);
    })($data["emptyFields"]) : null,
            complexField: isset($data["complexField"]) ? (function($input) {
    	/** @var array{uid?: string, nested?: mixed, array?: array<string>} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\DefaultsStructComplexField::fromArray($val);
    })($data["complexField"]) : null,
            partialComplexField: isset($data["partialComplexField"]) ? (function($input) {
    	/** @var array{uid?: string, intVal?: int} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\DefaultsStructPartialComplexField::fromArray($val);
    })($data["partialComplexField"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->allFields = $this->allFields;
        $data->partialFields = $this->partialFields;
        $data->emptyFields = $this->emptyFields;
        $data->complexField = $this->complexField;
        $data->partialComplexField = $this->partialComplexField;
        return $data;
    }
}
