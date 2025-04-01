<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class GridLayoutUsingValue implements \JsonSerializable
{
    public string $kind;

    public string $gridLayoutProperty;

    /**
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutKindType|null $kind
     * @param string|null $gridLayoutProperty
     */
    public function __construct(?\Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutKindType $kind = null, ?string $gridLayoutProperty = null)
    {
        $this->kind = $kind ?: \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutKindType;
        $this->gridLayoutProperty = $gridLayoutProperty ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{kind?: string, gridLayoutProperty?: string} $inputData */
        $data = $inputData;
        return new self(
            kind: isset($data["kind"]) ? /* ref to a non-struct, non-enum, this should have been inlined */ (function(array $input) { return $input; })($data["kind"]) : null,
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
