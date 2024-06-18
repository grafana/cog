<?php

namespace Grafana\Foundation\Defaults;

class SomeStruct implements \JsonSerializable {
    public bool $fieldBool;

    public string $fieldString;

    public string $fieldStringWithConstantValue;

    public float $fieldFloat32;

    public int $fieldInt32;

    /**
     * @param bool|null $fieldBool
     * @param string|null $fieldString
     * @param float|null $fieldFloat32
     * @param int|null $fieldInt32
     */
    public function __construct(?bool $fieldBool = null, ?string $fieldString = null, ?float $fieldFloat32 = null, ?int $fieldInt32 = null)
    {
        $this->fieldBool = $fieldBool ?: true;
        $this->fieldString = $fieldString ?: "foo";
        $this->fieldStringWithConstantValue = "auto";
    
        $this->fieldFloat32 = $fieldFloat32 ?: 42.42;
        $this->fieldInt32 = $fieldInt32 ?: 42;
    }

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
