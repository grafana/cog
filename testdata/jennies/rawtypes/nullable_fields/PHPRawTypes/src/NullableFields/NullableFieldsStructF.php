<?php

namespace Grafana\Foundation\NullableFields;

class NullableFieldsStructF implements \JsonSerializable
{
    public string $a;

    /**
     * @param string|null $a
     */
    public function __construct(?string $a = null)
    {
        $this->a = $a ?: "";
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->a !== $other->a) {
            return false;
        }
    
        return true;
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
