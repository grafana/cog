<?php

namespace Grafana\Foundation\Defaults;

class DefaultsStructComplexField implements \JsonSerializable
{
    public string $uid;

    public \Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested $nested;

    /**
     * @var array<string>
     */
    public array $array;

    /**
     * @param string|null $uid
     * @param \Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested|null $nested
     * @param array<string>|null $array
     */
    public function __construct(?string $uid = null, ?\Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested $nested = null, ?array $array = null)
    {
        $this->uid = $uid ?: "";
        $this->nested = $nested ?: new \Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested();
        $this->array = $array ?: [];
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{uid?: string, nested?: mixed, array?: array<string>} $inputData */
        $data = $inputData;
        return new self(
            uid: $data["uid"] ?? null,
            nested: isset($data["nested"]) ? (function($input) {
    	/** @var array{nestedVal?: string} */
    $val = $input;
    	return \Grafana\Foundation\Defaults\DefaultsStructComplexFieldNested::fromArray($val);
    })($data["nested"]) : null,
            array: $data["array"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->uid = $this->uid;
        $data->nested = $this->nested;
        $data->array = $this->array;
        return $data;
    }
}
