<?php

namespace Grafana\Foundation\DisjunctionAnonymous;

class StructA implements \JsonSerializable
{
    public string $field;

    /**
     * @param string|null $field
     */
    public function __construct(?string $field = null)
    {
        $this->field = $field ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{field?: string} $inputData */
        $data = $inputData;
        return new self(
            field: $data["field"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->field = $this->field;
        return $data;
    }
}
