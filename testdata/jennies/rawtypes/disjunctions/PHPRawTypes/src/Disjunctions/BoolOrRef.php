<?php

namespace Grafana\Foundation\Disjunctions;

class BoolOrRef implements \JsonSerializable
{
    public ?bool $bool;

    public ?\Grafana\Foundation\Disjunctions\SomeStruct $someStruct;

    /**
     * @param bool|null $bool
     * @param \Grafana\Foundation\Disjunctions\SomeStruct|null $someStruct
     */
    public function __construct(?bool $bool = null, ?\Grafana\Foundation\Disjunctions\SomeStruct $someStruct = null)
    {
        $this->bool = $bool;
        $this->someStruct = $someStruct;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{Bool?: bool, SomeStruct?: mixed} $inputData */
        $data = $inputData;
        return new self(
            bool: $data["Bool"] ?? null,
            someStruct: isset($data["SomeStruct"]) ? (function($input) {
    	/** @var array{Type?: string, FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\Disjunctions\SomeStruct::fromArray($val);
    })($data["SomeStruct"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->bool)) {
            $data->Bool = $this->bool;
        }
        if (isset($this->someStruct)) {
            $data->SomeStruct = $this->someStruct;
        }
        return $data;
    }
}
