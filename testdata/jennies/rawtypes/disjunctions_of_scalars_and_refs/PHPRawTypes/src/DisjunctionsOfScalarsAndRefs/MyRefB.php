<?php

namespace Grafana\Foundation\DisjunctionsOfScalarsAndRefs;

class MyRefB implements \JsonSerializable
{
    public int $bar;

    /**
     * @param int|null $bar
     */
    public function __construct(?int $bar = null)
    {
        $this->bar = $bar ?: 0;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
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
        /** @var array{bar?: int} $inputData */
        $data = $inputData;
        return new self(
            bar: $data["bar"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->bar = $this->bar;
        return $data;
    }
}
