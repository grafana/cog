<?php

namespace Grafana\Foundation\Types\Dashboard;

class FieldConfig implements \JsonSerializable {
    public ?string $unit;

    /**
     * @var mixed
     */
    public $custom;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
        ];
        if (isset($this->unit)) {
            $data["unit"] = $this->unit;
        }
        if (isset($this->custom)) {
            $data["custom"] = $this->custom;
        }
        return $data;
    }
}
