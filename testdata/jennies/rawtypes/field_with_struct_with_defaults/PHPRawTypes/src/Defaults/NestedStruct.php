<?php

namespace Grafana\Foundation\Defaults;

class NestedStruct implements \JsonSerializable {
    public string $stringVal;

    public int $intVal;

    /**
     * @param string|null $stringVal
     * @param int|null $intVal
     */
    public function __construct(?string $stringVal = null, ?int $intVal = null)
    {
        $this->stringVal = $stringVal ?: "";
        $this->intVal = $intVal ?: 0;
    }

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
