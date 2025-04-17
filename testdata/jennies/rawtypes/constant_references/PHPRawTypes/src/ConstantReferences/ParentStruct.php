<?php

namespace Grafana\Foundation\ConstantReferences;

class ParentStruct implements \JsonSerializable
{
    public \Grafana\Foundation\ConstantReferences\Enum $myEnum;

    /**
     * @param \Grafana\Foundation\ConstantReferences\Enum|null $myEnum
     */
    public function __construct(?\Grafana\Foundation\ConstantReferences\Enum $myEnum = null)
    {
        $this->myEnum = $myEnum ?: \Grafana\Foundation\ConstantReferences\Enum::ValueA();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{myEnum?: string} $inputData */
        $data = $inputData;
        return new self(
            myEnum: isset($data["myEnum"]) ? (function($input) { return \Grafana\Foundation\ConstantReferences\Enum::fromValue($input); })($data["myEnum"]) : null,
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
