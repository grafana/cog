<?php

namespace Grafana\Foundation\StructComplexFields;

class StringOrSomeOtherStruct implements \JsonSerializable
{
    public ?string $string;

    public ?\Grafana\Foundation\StructComplexFields\SomeOtherStruct $someOtherStruct;

    /**
     * @param string|null $string
     * @param \Grafana\Foundation\StructComplexFields\SomeOtherStruct|null $someOtherStruct
     */
    public function __construct(?string $string = null, ?\Grafana\Foundation\StructComplexFields\SomeOtherStruct $someOtherStruct = null)
    {
        $this->string = $string;
        $this->someOtherStruct = $someOtherStruct;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{String?: string, SomeOtherStruct?: mixed} $inputData */
        $data = $inputData;
        return new self(
            string: $data["String"] ?? null,
            someOtherStruct: isset($data["SomeOtherStruct"]) ? (function($input) {
    	/** @var array{FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\StructComplexFields\SomeOtherStruct::fromArray($val);
    })($data["SomeOtherStruct"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->string)) {
            $data->String = $this->string;
        }
        if (isset($this->someOtherStruct)) {
            $data->SomeOtherStruct = $this->someOtherStruct;
        }
        return $data;
    }
}
