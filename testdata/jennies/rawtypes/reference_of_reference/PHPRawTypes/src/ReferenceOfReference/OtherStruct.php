<?php

namespace Grafana\Foundation\ReferenceOfReference;

class OtherStruct extends \Grafana\Foundation\ReferenceOfReference\AnotherStruct
{
    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): static
    {
        $base = parent::fromArray($inputData);
        $obj = new self();
        foreach (get_object_vars($base) as $key => $value) {
            $obj->$key = $value;
        }
        return $obj;
    }
}
