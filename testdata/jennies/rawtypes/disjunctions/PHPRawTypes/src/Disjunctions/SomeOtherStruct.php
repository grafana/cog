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

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->type !== $other->type) {
            return false;
        }
    
        if ($this->foo !== $other->foo) {
            return false;
        }
    
        return true;
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
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->Type = $this->type;
        $data->Foo = $this->foo;
        return $data;
    }
}
