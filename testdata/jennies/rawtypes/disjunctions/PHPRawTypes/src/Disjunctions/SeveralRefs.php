<?php

namespace Grafana\Foundation\Disjunctions;

class SeveralRefs implements \JsonSerializable
{
    public ?\Grafana\Foundation\Disjunctions\SomeStruct $someStruct;

    public ?\Grafana\Foundation\Disjunctions\SomeOtherStruct $someOtherStruct;

    public ?\Grafana\Foundation\Disjunctions\YetAnotherStruct $yetAnotherStruct;

    /**
     * @param \Grafana\Foundation\Disjunctions\SomeStruct|null $someStruct
     * @param \Grafana\Foundation\Disjunctions\SomeOtherStruct|null $someOtherStruct
     * @param \Grafana\Foundation\Disjunctions\YetAnotherStruct|null $yetAnotherStruct
     */
    public function __construct(?\Grafana\Foundation\Disjunctions\SomeStruct $someStruct = null, ?\Grafana\Foundation\Disjunctions\SomeOtherStruct $someOtherStruct = null, ?\Grafana\Foundation\Disjunctions\YetAnotherStruct $yetAnotherStruct = null)
    {
        $this->someStruct = $someStruct;
        $this->someOtherStruct = $someOtherStruct;
        $this->yetAnotherStruct = $yetAnotherStruct;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{SomeStruct?: mixed, SomeOtherStruct?: mixed, YetAnotherStruct?: mixed} $inputData */
        $data = $inputData;
        return new self(
            someStruct: isset($data["SomeStruct"]) ? (function($input) {
    	/** @var array{Type?: string, FieldAny?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\Disjunctions\SomeStruct::fromArray($val);
    })($data["SomeStruct"]) : null,
            someOtherStruct: isset($data["SomeOtherStruct"]) ? (function($input) {
    	/** @var array{Type?: string, Foo?: string} */
    $val = $input;
    	return \Grafana\Foundation\Disjunctions\SomeOtherStruct::fromArray($val);
    })($data["SomeOtherStruct"]) : null,
            yetAnotherStruct: isset($data["YetAnotherStruct"]) ? (function($input) {
    	/** @var array{Type?: string, Bar?: int} */
    $val = $input;
    	return \Grafana\Foundation\Disjunctions\YetAnotherStruct::fromArray($val);
    })($data["YetAnotherStruct"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->someStruct)) {
            $data->SomeStruct = $this->someStruct;
        }
        if (isset($this->someOtherStruct)) {
            $data->SomeOtherStruct = $this->someOtherStruct;
        }
        if (isset($this->yetAnotherStruct)) {
            $data->YetAnotherStruct = $this->yetAnotherStruct;
        }
        return $data;
    }
}
