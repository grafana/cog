<?php

namespace Grafana\Foundation\Defaults;

class Struct implements \JsonSerializable {
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
