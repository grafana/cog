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

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->type !== $other->type) {
            return false;
        }
    
        if ($this->bar !== $other->bar) {
            return false;
        }
    
        return true;
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
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->Type = $this->type;
        $data->Bar = $this->bar;
        return $data;
    }
}
