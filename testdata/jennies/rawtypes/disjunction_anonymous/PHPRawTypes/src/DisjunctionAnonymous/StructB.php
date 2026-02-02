<?php

namespace Grafana\Foundation\DisjunctionAnonymous;

class StructB implements \JsonSerializable
{
    public int $type;

    /**
     * @param int|null $type
     */
    public function __construct(?int $type = null)
    {
        $this->type = $type ?: 0;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{type?: int} $inputData */
        $data = $inputData;
        return new self(
            type: $data["type"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->type = $this->type;
        return $data;
    }
}
