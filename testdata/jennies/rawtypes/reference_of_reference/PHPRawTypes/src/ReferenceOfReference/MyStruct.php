<?php

namespace Grafana\Foundation\ReferenceOfReference;

class MyStruct implements \JsonSerializable
{
    public ?\Grafana\Foundation\ReferenceOfReference\OtherStruct $field;

    /**
     * @param \Grafana\Foundation\ReferenceOfReference\OtherStruct|null $field
     */
    public function __construct(?\Grafana\Foundation\ReferenceOfReference\OtherStruct $field = null)
    {
        $this->field = $field;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if (($this->field === null) !== ($other->field === null)) {
            return false;
        }
        if ($this->field !== null) {
            if (!$this->field->equals($other->field)) {
                return false;
            }
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{field?: mixed} $inputData */
        $data = $inputData;
        return new self(
            field: isset($data["field"]) ? (function($input) {
    	/** @var array{a?: string} */
    $val = $input;
    	return \Grafana\Foundation\ReferenceOfReference\OtherStruct::fromArray($val);
    })($data["field"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->field)) {
            $data->field = $this->field;
        }
        return $data;
    }
}
