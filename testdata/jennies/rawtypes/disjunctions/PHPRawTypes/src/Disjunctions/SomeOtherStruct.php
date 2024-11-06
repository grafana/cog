<?php

namespace Grafana\Foundation\Disjunctions;

class SomeOtherStruct implements \JsonSerializable
{
    public string $type;

    public string $foo;

    /**
     * @param string|null $foo
     */
    public function __construct(?string $foo = null)
    {
        $this->type = "some-other-struct";
    
        $this->foo = $foo ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{Type?: string, Foo?: string} $inputData */
        $data = $inputData;
        return new self(
            foo: $data["Foo"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "Type" => $this->type,
            "Foo" => $this->foo,
        ];
        return $data;
    }
}
