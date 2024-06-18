<?php

namespace Grafana\Foundation\Dashboard;

class FieldConfig implements \JsonSerializable {
    public ?string $unit;

    /**
     * @var mixed|null
     */
    public $custom;

    /**
     * @param string|null $unit
     * @param mixed|null $custom
     */
    public function __construct(?string $unit = null,  $custom = null)
    {
        $this->unit = $unit;
        $this->custom = $custom;
    }

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
