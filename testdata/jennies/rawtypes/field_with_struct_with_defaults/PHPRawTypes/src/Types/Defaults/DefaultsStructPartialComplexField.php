<?php

namespace Grafana\Foundation\Types\Defaults;

class DefaultsStructPartialComplexField implements \JsonSerializable {
    public string $uid;

    public int $intVal;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "uid" => $this->uid,
            "intVal" => $this->intVal,
        ];
        return $data;
    }
}
