<?php

namespace Grafana\Foundation\Types\StructOptionalFields;

class SomeOtherStruct implements \JsonSerializable {
    /**
     * @var mixed
     */
    public $fieldAny;

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
