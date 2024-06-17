<?php

namespace Grafana\Foundation\Types\Dashboard;

class FieldConfigSource implements \JsonSerializable {
    public ?\Grafana\Foundation\Types\Dashboard\FieldConfig $defaults;

    /**
     * @param \Grafana\Foundation\Types\Dashboard\FieldConfig|null $defaults
     */
    public function __construct(?\Grafana\Foundation\Types\Dashboard\FieldConfig $defaults = null)
    {
        $this->defaults = $defaults;
    }

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
