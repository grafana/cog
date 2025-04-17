<?php

namespace Grafana\Foundation\Basic;

/**
 * This
 * is
 * a
 * comment
 */
class SomeStruct implements \JsonSerializable
{
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
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldAny?: mixed, FieldBool?: bool, FieldBytes?: string, FieldString?: string, FieldStringWithConstantValue?: string, FieldFloat32?: float, FieldFloat64?: float, FieldUint8?: int, FieldUint16?: int, FieldUint32?: int, FieldUint64?: int, FieldInt8?: int, FieldInt16?: int, FieldInt32?: int, FieldInt64?: int} $inputData */
        $data = $inputData;
        return new self(
            fieldAny: $data["FieldAny"] ?? null,
            fieldBool: $data["FieldBool"] ?? null,
            fieldBytes: $data["FieldBytes"] ?? null,
            fieldString: $data["FieldString"] ?? null,
            fieldFloat32: $data["FieldFloat32"] ?? null,
            fieldFloat64: $data["FieldFloat64"] ?? null,
            fieldUint8: $data["FieldUint8"] ?? null,
            fieldUint16: $data["FieldUint16"] ?? null,
            fieldUint32: $data["FieldUint32"] ?? null,
            fieldUint64: $data["FieldUint64"] ?? null,
            fieldInt8: $data["FieldInt8"] ?? null,
            fieldInt16: $data["FieldInt16"] ?? null,
            fieldInt32: $data["FieldInt32"] ?? null,
            fieldInt64: $data["FieldInt64"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->FieldAny = $this->fieldAny;
        $data->FieldBool = $this->fieldBool;
        $data->FieldBytes = $this->fieldBytes;
        $data->FieldString = $this->fieldString;
        $data->FieldStringWithConstantValue = $this->fieldStringWithConstantValue;
        $data->FieldFloat32 = $this->fieldFloat32;
        $data->FieldFloat64 = $this->fieldFloat64;
        $data->FieldUint8 = $this->fieldUint8;
        $data->FieldUint16 = $this->fieldUint16;
        $data->FieldUint32 = $this->fieldUint32;
        $data->FieldUint64 = $this->fieldUint64;
        $data->FieldInt8 = $this->fieldInt8;
        $data->FieldInt16 = $this->fieldInt16;
        $data->FieldInt32 = $this->fieldInt32;
        $data->FieldInt64 = $this->fieldInt64;
        return $data;
    }
}
