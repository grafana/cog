<?php

namespace Grafana\Foundation\Types\Basic;

/**
 * This
 * is
 * a
 * comment
 */
class SomeStruct implements \JsonSerializable {
    /**
     * Anything can go in there.
     * Really, anything.
     * @var mixed
     */
    public $fieldAny;

    public bool $fieldBool;

    public string $fieldBytes;

    public string $fieldString;

    public string $fieldStringWithConstantValue;

    public float $fieldFloat32;

    public float $fieldFloat64;

    public int $fieldUint8;

    public int $fieldUint16;

    public int $fieldUint32;

    public int $fieldUint64;

    public int $fieldInt8;

    public int $fieldInt16;

    public int $fieldInt32;

    public int $fieldInt64;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "FieldAny" => $this->fieldAny,
            "FieldBool" => $this->fieldBool,
            "FieldBytes" => $this->fieldBytes,
            "FieldString" => $this->fieldString,
            "FieldStringWithConstantValue" => $this->fieldStringWithConstantValue,
            "FieldFloat32" => $this->fieldFloat32,
            "FieldFloat64" => $this->fieldFloat64,
            "FieldUint8" => $this->fieldUint8,
            "FieldUint16" => $this->fieldUint16,
            "FieldUint32" => $this->fieldUint32,
            "FieldUint64" => $this->fieldUint64,
            "FieldInt8" => $this->fieldInt8,
            "FieldInt16" => $this->fieldInt16,
            "FieldInt32" => $this->fieldInt32,
            "FieldInt64" => $this->fieldInt64,
        ];
        return $data;
    }
}
