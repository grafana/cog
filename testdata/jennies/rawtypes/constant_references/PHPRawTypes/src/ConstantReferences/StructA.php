<?php

namespace Grafana\Foundation\ConstantReferences;

class StructA implements \JsonSerializable
{
    public \Grafana\Foundation\ConstantReferences\Enum $myEnum;

    public function __construct()
    {
        $this->myEnum = \Grafana\Foundation\ConstantReferences\Enum::valueA();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        return new self(
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->myEnum = $this->myEnum;
        return $data;
    }
}
