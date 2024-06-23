<?php

namespace Grafana\Foundation\StructComplexFields;

class SomeOtherStruct implements \JsonSerializable
{
    /**
     * @var mixed
     */
    public $fieldAny;

    /**
     * @param mixed|null $fieldAny
     */
    public function __construct( $fieldAny = null)
    {
        $this->fieldAny = $fieldAny ?: null;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldAny?: mixed} $inputData */
        $data = $inputData;
        return new self(
            fieldAny: $data["FieldAny"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "FieldAny" => $this->fieldAny,
        ];
        return $data;
    }
}
