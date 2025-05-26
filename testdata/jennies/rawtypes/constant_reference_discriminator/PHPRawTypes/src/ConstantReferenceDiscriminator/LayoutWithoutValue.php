<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class LayoutWithoutValue implements \JsonSerializable
{
    public ?\Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutWithoutValue $gridLayoutWithoutValue;

    public ?\Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutWithoutValue $rowsLayoutWithoutValue;

    /**
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutWithoutValue|null $gridLayoutWithoutValue
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutWithoutValue|null $rowsLayoutWithoutValue
     */
    public function __construct(?\Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutWithoutValue $gridLayoutWithoutValue = null, ?\Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutWithoutValue $rowsLayoutWithoutValue = null)
    {
        $this->gridLayoutWithoutValue = $gridLayoutWithoutValue;
        $this->rowsLayoutWithoutValue = $rowsLayoutWithoutValue;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{GridLayoutWithoutValue?: mixed, RowsLayoutWithoutValue?: mixed} $inputData */
        $data = $inputData;
        return new self(
            gridLayoutWithoutValue: isset($data["GridLayoutWithoutValue"]) ? (function($input) {
    	/** @var array{kind?: "GridLayout", gridLayoutProperty?: string} */
    $val = $input;
    	return \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutWithoutValue::fromArray($val);
    })($data["GridLayoutWithoutValue"]) : null,
            rowsLayoutWithoutValue: isset($data["RowsLayoutWithoutValue"]) ? (function($input) {
    	/** @var array{kind?: "RowsLayout", rowsLayoutProperty?: string} */
    $val = $input;
    	return \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutWithoutValue::fromArray($val);
    })($data["RowsLayoutWithoutValue"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->gridLayoutWithoutValue)) {
            $data->GridLayoutWithoutValue = $this->gridLayoutWithoutValue;
        }
        if (isset($this->rowsLayoutWithoutValue)) {
            $data->RowsLayoutWithoutValue = $this->rowsLayoutWithoutValue;
        }
        return $data;
    }
}
