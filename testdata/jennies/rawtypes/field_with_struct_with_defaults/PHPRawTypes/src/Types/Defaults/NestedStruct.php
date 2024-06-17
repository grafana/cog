<?php

namespace Grafana\Foundation\Types\Defaults;

class NestedStruct implements \JsonSerializable {
    public string $stringVal;

    public int $intVal;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "stringVal" => $this->stringVal,
            "intVal" => $this->intVal,
        ];
        return $data;
    }
}
