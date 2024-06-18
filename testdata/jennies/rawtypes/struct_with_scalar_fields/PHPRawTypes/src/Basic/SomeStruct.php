<?php

namespace Grafana\Foundation\Basic;

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
     * @param mixed|null $fieldAny
     * @param bool|null $fieldBool
     * @param string|null $fieldBytes
     * @param string|null $fieldString
     * @param float|null $fieldFloat32
     * @param float|null $fieldFloat64
     * @param int|null $fieldUint8
     * @param int|null $fieldUint16
     * @param int|null $fieldUint32
     * @param int|null $fieldUint64
     * @param int|null $fieldInt8
     * @param int|null $fieldInt16
     * @param int|null $fieldInt32
     * @param int|null $fieldInt64
     */
    public function __construct( $fieldAny = null, ?bool $fieldBool = null, ?string $fieldBytes = null, ?string $fieldString = null, ?float $fieldFloat32 = null, ?float $fieldFloat64 = null, ?int $fieldUint8 = null, ?int $fieldUint16 = null, ?int $fieldUint32 = null, ?int $fieldUint64 = null, ?int $fieldInt8 = null, ?int $fieldInt16 = null, ?int $fieldInt32 = null, ?int $fieldInt64 = null)
    {
        $this->fieldAny = $fieldAny ?: null;
        $this->fieldBool = $fieldBool ?: false;
        $this->fieldBytes = $fieldBytes ?: "";
        $this->fieldString = $fieldString ?: "";
        $this->fieldStringWithConstantValue = "auto";
    
        $this->fieldFloat32 = $fieldFloat32 ?: 0;
        $this->fieldFloat64 = $fieldFloat64 ?: 0;
        $this->fieldUint8 = $fieldUint8 ?: 0;
        $this->fieldUint16 = $fieldUint16 ?: 0;
        $this->fieldUint32 = $fieldUint32 ?: 0;
        $this->fieldUint64 = $fieldUint64 ?: 0;
        $this->fieldInt8 = $fieldInt8 ?: 0;
        $this->fieldInt16 = $fieldInt16 ?: 0;
        $this->fieldInt32 = $fieldInt32 ?: 0;
        $this->fieldInt64 = $fieldInt64 ?: 0;
    }

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
