<?php

namespace Grafana\Foundation\Defaults;

class SomeStruct implements \JsonSerializable
{
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
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{fieldBool?: bool, fieldString?: string, FieldStringWithConstantValue?: string, FieldFloat32?: float, FieldInt32?: int} $inputData */
        $data = $inputData;
        return new self(
            fieldBool: $data["fieldBool"] ?? null,
            fieldString: $data["fieldString"] ?? null,
            fieldFloat32: $data["FieldFloat32"] ?? null,
            fieldInt32: $data["FieldInt32"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->fieldBool = $this->fieldBool;
        $data->fieldString = $this->fieldString;
        $data->FieldStringWithConstantValue = $this->fieldStringWithConstantValue;
        $data->FieldFloat32 = $this->fieldFloat32;
        $data->FieldInt32 = $this->fieldInt32;
        return $data;
    }
}
