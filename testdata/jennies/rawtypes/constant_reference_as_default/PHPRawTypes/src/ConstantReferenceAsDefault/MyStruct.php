<?php

namespace Grafana\Foundation\ConstantReferenceAsDefault;

class MyStruct implements \JsonSerializable
{
    public string $aString;

    public ?string $optString;

    public function __construct()
    {
        $this->aString = \Grafana\Foundation\ConstantReferenceAsDefault\ConstantRefString;
        $this->optString = \Grafana\Foundation\ConstantReferenceAsDefault\ConstantRefString;
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
        $data->aString = $this->aString;
        if (isset($this->optString)) {
            $data->optString = $this->optString;
        }
        return $data;
    }
}
