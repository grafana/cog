<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class RowsLayoutWithoutValue implements \JsonSerializable
{
    public string $kind;

    public string $rowsLayoutProperty;

    /**
     * @param string|null $rowsLayoutProperty
     */
    public function __construct(?string $rowsLayoutProperty = null)
    {
        $this->kind = \Grafana\Foundation\ConstantReferenceDiscriminator\Constants::ROWS_LAYOUT_KIND_TYPE;
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
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->kind = $this->kind;
        $data->rowsLayoutProperty = $this->rowsLayoutProperty;
        return $data;
    }
}
