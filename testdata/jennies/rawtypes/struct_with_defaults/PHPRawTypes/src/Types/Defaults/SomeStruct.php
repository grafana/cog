<?php

namespace Grafana\Foundation\Types\Defaults;

class SomeStruct implements \JsonSerializable {
    public bool $fieldBool;

    public string $fieldString;

    public string $fieldStringWithConstantValue;

    public float $fieldFloat32;

    public int $fieldInt32;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "fieldBool" => $this->fieldBool,
            "fieldString" => $this->fieldString,
            "FieldStringWithConstantValue" => $this->fieldStringWithConstantValue,
            "FieldFloat32" => $this->fieldFloat32,
            "FieldInt32" => $this->fieldInt32,
        ];
        return $data;
    }
}
