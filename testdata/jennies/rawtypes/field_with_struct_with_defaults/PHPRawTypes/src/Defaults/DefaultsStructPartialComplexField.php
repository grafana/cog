<?php

namespace Grafana\Foundation\Defaults;

class DefaultsStructPartialComplexField implements \JsonSerializable {
    public string $uid;

    public int $intVal;

    /**
     * @param string|null $uid
     * @param int|null $intVal
     */
    public function __construct(?string $uid = null, ?int $intVal = null)
    {
        $this->uid = $uid ?: "";
        $this->intVal = $intVal ?: 0;
    }

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
