<?php

namespace Grafana\Foundation\DefaultDisjunctionValue;

class ValueC implements \JsonSerializable
{
    public string $type;

    public float $other;

    /**
     * @param float|null $other
     */
    public function __construct(?float $other = null)
    {
        $this->type = "C";
    
        $this->other = $other ?: 0;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{type?: string, other?: float} $inputData */
        $data = $inputData;
        return new self(
            other: $data["other"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "type" => $this->type,
            "other" => $this->other,
        ];
        return $data;
    }
}
