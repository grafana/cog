<?php

namespace Grafana\Foundation\Basic;

/**
 * @deprecated This object is deprecated, use NewStruct instead.
 */
class SomeStruct implements \JsonSerializable
{
    public string $fieldString;

    /**
     * @param string|null $fieldString
     */
    public function __construct(?string $fieldString = null)
    {
        $this->fieldString = $fieldString ?: "";
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->fieldString !== $other->fieldString) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldString?: string} $inputData */
        $data = $inputData;
        return new self(
            fieldString: $data["FieldString"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->FieldString = $this->fieldString;
        return $data;
    }
}
