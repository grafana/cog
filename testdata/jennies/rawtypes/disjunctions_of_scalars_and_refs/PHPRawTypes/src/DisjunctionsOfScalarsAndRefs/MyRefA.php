<?php

namespace Grafana\Foundation\DisjunctionsOfScalarsAndRefs;

class MyRefA implements \JsonSerializable
{
    public string $foo;

    /**
     * @param string|null $foo
     */
    public function __construct(?string $foo = null)
    {
        $this->foo = $foo ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{foo?: string} $inputData */
        $data = $inputData;
        return new self(
            foo: $data["foo"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->foo = $this->foo;
        return $data;
    }
}
