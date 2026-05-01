<?php

namespace Grafana\Foundation\DisjunctionsOfRefsWithoutDiscriminator;

class TypeB implements \JsonSerializable
{
    public int $fieldB;

    /**
     * @param int|null $fieldB
     */
    public function __construct(?int $fieldB = null)
    {
        $this->fieldB = $fieldB ?: 0;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->fieldB !== $other->fieldB) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{fieldB?: int} $inputData */
        $data = $inputData;
        return new self(
            fieldB: $data["fieldB"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->fieldB = $this->fieldB;
        return $data;
    }
}
