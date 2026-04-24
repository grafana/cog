<?php

namespace Grafana\Foundation\ReferenceOfReference;

class AnotherStruct implements \JsonSerializable
{
    public string $a;

    /**
     * @param string|null $a
     */
    public function __construct(?string $a = null)
    {
        $this->a = $a ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{a?: string} $inputData */
        $data = $inputData;
        return new self(
            a: $data["a"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->a = $this->a;
        return $data;
    }
}
