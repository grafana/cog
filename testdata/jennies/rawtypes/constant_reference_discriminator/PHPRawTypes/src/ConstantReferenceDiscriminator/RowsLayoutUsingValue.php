<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class RowsLayoutUsingValue implements \JsonSerializable
{
    public string $kind;

    public string $rowsLayoutProperty;

    /**
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutKindType|null $kind
     * @param string|null $rowsLayoutProperty
     */
    public function __construct(?\Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutKindType $kind = null, ?string $rowsLayoutProperty = null)
    {
        $this->kind = $kind ?: \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutKindType;
        $this->rowsLayoutProperty = $rowsLayoutProperty ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{kind?: string, rowsLayoutProperty?: string} $inputData */
        $data = $inputData;
        return new self(
            kind: isset($data["kind"]) ? /* ref to a non-struct, non-enum, this should have been inlined */ (function(array $input) { return $input; })($data["kind"]) : null,
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
