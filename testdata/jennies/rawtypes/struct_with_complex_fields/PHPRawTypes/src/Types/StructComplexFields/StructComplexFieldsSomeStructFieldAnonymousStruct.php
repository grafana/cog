<?php

namespace Grafana\Foundation\Types\StructComplexFields;

class StructComplexFieldsSomeStructFieldAnonymousStruct implements \JsonSerializable {
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

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "FieldAny" => $this->fieldAny,
        ];
        return $data;
    }
}
