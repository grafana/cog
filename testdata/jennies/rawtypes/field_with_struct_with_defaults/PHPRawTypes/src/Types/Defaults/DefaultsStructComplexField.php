<?php

namespace Grafana\Foundation\Types\Defaults;

class DefaultsStructComplexField implements \JsonSerializable {
    public string $uid;

    public \Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested $nested;

    /**
     * @var array<string>
     */
    public array $array;

    /**
     * @param string|null $uid
     * @param \Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested|null $nested
     * @param array<string>|null $array
     */
    public function __construct(?string $uid = null, ?\Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested $nested = null, ?array $array = null)
    {
        $this->uid = $uid ?: "";
        $this->nested = $nested ?: new \Grafana\Foundation\Types\Defaults\DefaultsStructComplexFieldNested();
        $this->array = $array ?: [];
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "uid" => $this->uid,
            "nested" => $this->nested,
            "array" => $this->array,
        ];
        return $data;
    }
}
