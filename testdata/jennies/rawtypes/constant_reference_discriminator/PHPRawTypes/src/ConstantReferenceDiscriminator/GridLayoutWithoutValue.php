<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class GridLayoutWithoutValue implements \JsonSerializable
{
    public string $kind;

    public string $gridLayoutProperty;

    /**
     * @param string|null $gridLayoutProperty
     */
    public function __construct(?string $gridLayoutProperty = null)
    {
        $this->kind = \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutKindType;
        $this->gridLayoutProperty = $gridLayoutProperty ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{kind?: "GridLayout", gridLayoutProperty?: string} $inputData */
        $data = $inputData;
        return new self(
            gridLayoutProperty: $data["gridLayoutProperty"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "kind" => $this->kind,
            "gridLayoutProperty" => $this->gridLayoutProperty,
        ];
        return $data;
    }
}
