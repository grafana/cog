<?php

namespace Grafana\Foundation\DisjunctionsOfRefsWithoutDiscriminator;

class TypeA implements \JsonSerializable
{
    public string $fieldA;

    /**
     * @param string|null $fieldA
     */
    public function __construct(?string $fieldA = null)
    {
        $this->fieldA = $fieldA ?: "";
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->fieldA !== $other->fieldA) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{fieldA?: string} $inputData */
        $data = $inputData;
        return new self(
            fieldA: $data["fieldA"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->fieldA = $this->fieldA;
        return $data;
    }
}
