<?php

namespace Grafana\Foundation\Types\Dashboard;

class FieldConfigSource implements \JsonSerializable {
    public ?\Grafana\Foundation\Types\Dashboard\FieldConfig $defaults;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
        ];
        if (isset($this->defaults)) {
            $data["defaults"] = $this->defaults;
        }
        return $data;
    }
}
