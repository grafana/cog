<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class RowsLayoutUsingValue implements \JsonSerializable
{
    public string $kind;

    public string $rowsLayoutProperty;

    /**
     * @param string|null $rowsLayoutProperty
     */
    public function __construct(?string $rowsLayoutProperty = null)
    {
        $this->kind = \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutKindType;
        $this->rowsLayoutProperty = $rowsLayoutProperty ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{kind?: "RowsLayout", rowsLayoutProperty?: string} $inputData */
        $data = $inputData;
        return new self(
            rowsLayoutProperty: $data["rowsLayoutProperty"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "kind" => $this->kind,
            "rowsLayoutProperty" => $this->rowsLayoutProperty,
        ];
        return $data;
    }
}
