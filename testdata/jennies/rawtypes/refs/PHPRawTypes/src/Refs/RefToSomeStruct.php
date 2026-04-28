<?php

namespace Grafana\Foundation\Refs;

class RefToSomeStruct extends \Grafana\Foundation\Refs\SomeStruct
{
    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): static
    {
        $base = parent::fromArray($inputData);
        $obj = new static();
        foreach (get_object_vars($base) as $key => $value) {
            $obj->$key = $value;
        }
        return $obj;
    }
}
