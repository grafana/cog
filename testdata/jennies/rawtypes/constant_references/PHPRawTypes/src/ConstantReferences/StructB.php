<?php

namespace Grafana\Foundation\ConstantReferences;

class StructB implements \JsonSerializable
{
    public \Grafana\Foundation\ConstantReferences\Enum $myEnum;

    public string $myValue;

    /**
     * @param string|null $myValue
     */
    public function __construct(?string $myValue = null)
    {
        $this->myEnum = \Grafana\Foundation\ConstantReferences\Enum::valueA();
        $this->myValue = $myValue ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{myEnum?: "ValueB", myValue?: string} $inputData */
        $data = $inputData;
        return new self(
            myValue: $data["myValue"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->myEnum = $this->myEnum;
        $data->myValue = $this->myValue;
        return $data;
    }
}
