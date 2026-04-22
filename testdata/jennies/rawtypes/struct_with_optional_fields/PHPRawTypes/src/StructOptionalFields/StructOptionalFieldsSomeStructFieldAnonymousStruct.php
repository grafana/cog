<?php

namespace Grafana\Foundation\StructOptionalFields;

class StructOptionalFieldsSomeStructFieldAnonymousStruct implements \JsonSerializable
{
    /**
     * @var mixed
     */
    public $fieldAny;

    /**
     * @param mixed|null $fieldAny
     */
    public function __construct( $fieldAny = null)
    {
        $this->fieldAny = $fieldAny ?: null;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->fieldAny !== $other->fieldAny) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{FieldAny?: mixed} $inputData */
        $data = $inputData;
        return new self(
            fieldAny: $data["FieldAny"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->FieldAny = $this->fieldAny;
        return $data;
    }
}
