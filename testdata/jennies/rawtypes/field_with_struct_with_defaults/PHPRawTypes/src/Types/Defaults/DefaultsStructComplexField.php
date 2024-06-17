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
