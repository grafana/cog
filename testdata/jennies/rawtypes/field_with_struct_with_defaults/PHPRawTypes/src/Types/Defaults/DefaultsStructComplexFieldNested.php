<?php

namespace Grafana\Foundation\Types\Defaults;

class DefaultsStructComplexFieldNested implements \JsonSerializable {
    public string $nestedVal;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "nestedVal" => $this->nestedVal,
        ];
        return $data;
    }
}
