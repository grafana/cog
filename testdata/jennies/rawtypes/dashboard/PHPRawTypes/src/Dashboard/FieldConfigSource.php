<?php

namespace Grafana\Foundation\Dashboard;

class FieldConfigSource implements \JsonSerializable {
    public ?\Grafana\Foundation\Dashboard\FieldConfig $defaults;

    /**
     * @param \Grafana\Foundation\Dashboard\FieldConfig|null $defaults
     */
    public function __construct(?\Grafana\Foundation\Dashboard\FieldConfig $defaults = null)
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
