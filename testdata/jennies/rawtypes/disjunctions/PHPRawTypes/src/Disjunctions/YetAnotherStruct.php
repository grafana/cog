<?php

namespace Grafana\Foundation\Disjunctions;

class YetAnotherStruct implements \JsonSerializable
{
    public string $type;

    public int $bar;

    /**
     * @param int|null $bar
     */
    public function __construct(?int $bar = null)
    {
        $this->type = "yet-another-struct";
    
        $this->bar = $bar ?: 0;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{Type?: string, Bar?: int} $inputData */
        $data = $inputData;
        return new self(
            bar: $data["Bar"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "Type" => $this->type,
            "Bar" => $this->bar,
        ];
        return $data;
    }
}