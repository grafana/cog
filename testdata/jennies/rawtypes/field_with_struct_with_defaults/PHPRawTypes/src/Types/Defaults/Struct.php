<?php

namespace Grafana\Foundation\Types\Defaults;

class Struct implements \JsonSerializable {
    public \Grafana\Foundation\Types\Defaults\NestedStruct $allFields;

    public \Grafana\Foundation\Types\Defaults\NestedStruct $partialFields;

    public \Grafana\Foundation\Types\Defaults\NestedStruct $emptyFields;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructComplexField $complexField;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructPartialComplexField $partialComplexField;

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
