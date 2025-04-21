<?php

namespace Grafana\Foundation\Dashboard;

class FieldConfigSource implements \JsonSerializable
{
    public ?\Grafana\Foundation\Dashboard\FieldConfig $defaults;

    /**
     * @param \Grafana\Foundation\Dashboard\FieldConfig|null $defaults
     */
    public function __construct(?\Grafana\Foundation\Dashboard\FieldConfig $defaults = null)
    {
        $this->defaults = $defaults;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{defaults?: mixed} $inputData */
        $data = $inputData;
        return new self(
            defaults: isset($data["defaults"]) ? (function($input) {
    	/** @var array{unit?: string, custom?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\Dashboard\FieldConfig::fromArray($val);
    })($data["defaults"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->defaults)) {
            $data->defaults = $this->defaults;
        }
        return $data;
    }
}
