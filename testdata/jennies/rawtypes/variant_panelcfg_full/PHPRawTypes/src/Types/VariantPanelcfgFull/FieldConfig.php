<?php

namespace Grafana\Foundation\Types\VariantPanelcfgFull;

class FieldConfig implements \JsonSerializable {
    public string $timeseriesFieldConfigOption;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "timeseries_field_config_option" => $this->timeseriesFieldConfigOption,
        ];
        return $data;
    }
}
