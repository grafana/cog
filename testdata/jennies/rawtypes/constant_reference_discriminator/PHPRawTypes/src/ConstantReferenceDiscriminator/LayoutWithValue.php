<?php

namespace Grafana\Foundation\ConstantReferenceDiscriminator;

class LayoutWithValue implements \JsonSerializable
{
    public ?\Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutUsingValue $gridLayoutUsingValue;

    public ?\Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutUsingValue $rowsLayoutUsingValue;

    /**
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutUsingValue|null $gridLayoutUsingValue
     * @param \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutUsingValue|null $rowsLayoutUsingValue
     */
    public function __construct(?\Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutUsingValue $gridLayoutUsingValue = null, ?\Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutUsingValue $rowsLayoutUsingValue = null)
    {
        $this->gridLayoutUsingValue = $gridLayoutUsingValue;
        $this->rowsLayoutUsingValue = $rowsLayoutUsingValue;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{GridLayoutUsingValue?: mixed, RowsLayoutUsingValue?: mixed} $inputData */
        $data = $inputData;
        return new self(
            gridLayoutUsingValue: isset($data["GridLayoutUsingValue"]) ? (function($input) {
    	/** @var array{kind?: "GridLayout", gridLayoutProperty?: string} */
    $val = $input;
    	return \Grafana\Foundation\ConstantReferenceDiscriminator\GridLayoutUsingValue::fromArray($val);
    })($data["GridLayoutUsingValue"]) : null,
            rowsLayoutUsingValue: isset($data["RowsLayoutUsingValue"]) ? (function($input) {
    	/** @var array{kind?: "RowsLayout", rowsLayoutProperty?: string} */
    $val = $input;
    	return \Grafana\Foundation\ConstantReferenceDiscriminator\RowsLayoutUsingValue::fromArray($val);
    })($data["RowsLayoutUsingValue"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->gridLayoutUsingValue)) {
            $data->GridLayoutUsingValue = $this->gridLayoutUsingValue;
        }
        if (isset($this->rowsLayoutUsingValue)) {
            $data->RowsLayoutUsingValue = $this->rowsLayoutUsingValue;
        }
        return $data;
    }
}
