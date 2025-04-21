<?php

namespace Grafana\Foundation\Disjunctions;

class SomeStruct implements \JsonSerializable
{
    public string $type;

    /**
     * @var mixed
     */
    public $fieldAny;

    /**
     * @param mixed|null $fieldAny
     */
    public function __construct( $fieldAny = null)
    {
        $this->type = "some-struct";
    
        $this->fieldAny = $fieldAny ?: null;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{Type?: string, FieldAny?: mixed} $inputData */
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
        $data->Type = $this->type;
        $data->FieldAny = $this->fieldAny;
        return $data;
    }
}
