<?php

namespace Grafana\Foundation\ConstantReferences;

class Struct implements \JsonSerializable
{
    public string $myValue;

    public \Grafana\Foundation\ConstantReferences\Enum $myEnum;

    /**
     * @param string|null $myValue
     * @param \Grafana\Foundation\ConstantReferences\Enum|null $myEnum
     */
    public function __construct(?string $myValue = null, ?\Grafana\Foundation\ConstantReferences\Enum $myEnum = null)
    {
        $this->myValue = $myValue ?: "";
        $this->myEnum = $myEnum ?: \Grafana\Foundation\ConstantReferences\Enum::ValueA();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{myValue?: string, myEnum?: string} $inputData */
        $data = $inputData;
        return new self(
            myValue: $data["myValue"] ?? null,
            myEnum: isset($data["myEnum"]) ? (function($input) { return \Grafana\Foundation\ConstantReferences\Enum::fromValue($input); })($data["myEnum"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->myValue = $this->myValue;
        $data->myEnum = $this->myEnum;
        return $data;
    }
}
