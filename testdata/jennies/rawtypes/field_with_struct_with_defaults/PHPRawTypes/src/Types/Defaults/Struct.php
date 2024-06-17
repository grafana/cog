<?php

namespace Grafana\Foundation\Types\Defaults;

class Struct implements \JsonSerializable {
    public \Grafana\Foundation\Types\Defaults\NestedStruct $allFields;

    public \Grafana\Foundation\Types\Defaults\NestedStruct $partialFields;

    public \Grafana\Foundation\Types\Defaults\NestedStruct $emptyFields;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructComplexField $complexField;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructPartialComplexField $partialComplexField;

    /**
     * @param \Grafana\Foundation\Types\Defaults\NestedStruct|null $allFields
     * @param \Grafana\Foundation\Types\Defaults\NestedStruct|null $partialFields
     * @param \Grafana\Foundation\Types\Defaults\NestedStruct|null $emptyFields
     * @param \Grafana\Foundation\Types\Defaults\DefaultsStructComplexField|null $complexField
     * @param \Grafana\Foundation\Types\Defaults\DefaultsStructPartialComplexField|null $partialComplexField
     */
    public function __construct(?\Grafana\Foundation\Types\Defaults\NestedStruct $allFields = null, ?\Grafana\Foundation\Types\Defaults\NestedStruct $partialFields = null, ?\Grafana\Foundation\Types\Defaults\NestedStruct $emptyFields = null, ?\Grafana\Foundation\Types\Defaults\DefaultsStructComplexField $complexField = null, ?\Grafana\Foundation\Types\Defaults\DefaultsStructPartialComplexField $partialComplexField = null)
    {
        $this->allFields = $allFields ?: new \Grafana\Foundation\Types\Defaults\NestedStruct(intVal: 3, stringVal: "hello");
        $this->partialFields = $partialFields ?: new \Grafana\Foundation\Types\Defaults\NestedStruct(intVal: 3);
        $this->emptyFields = $emptyFields ?: new \Grafana\Foundation\Types\Defaults\NestedStruct();
        $this->complexField = $complexField ?: new \Grafana\Foundation\Types\Defaults\DefaultsStructComplexField(array: ["hello"], nested: new \Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested(nestedVal: "nested"), uid: "myUID");
        $this->partialComplexField = $partialComplexField ?: new \Grafana\Foundation\Types\Defaults\DefaultsStructPartialComplexField();
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "allFields" => $this->allFields,
            "partialFields" => $this->partialFields,
            "emptyFields" => $this->emptyFields,
            "complexField" => $this->complexField,
            "partialComplexField" => $this->partialComplexField,
        ];
        return $data;
    }
}
