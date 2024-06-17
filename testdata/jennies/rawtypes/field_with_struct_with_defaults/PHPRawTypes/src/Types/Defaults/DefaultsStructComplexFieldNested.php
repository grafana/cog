<?php

namespace Grafana\Foundation\Types\Defaults;

class DefaultsStructComplexFieldNested implements \JsonSerializable {
    public string $nestedVal;

    /**
     * @param string|null $nestedVal
     */
    public function __construct(?string $nestedVal = null)
    {
        $this->nestedVal = $nestedVal ?: "";
    }

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
