<?php

namespace Grafana\Foundation\VariantPanelcfgFull;

class FieldConfig implements \JsonSerializable {
    public string $timeseriesFieldConfigOption;

    /**
     * @param string|null $timeseriesFieldConfigOption
     */
    public function __construct(?string $timeseriesFieldConfigOption = null)
    {
        $this->timeseriesFieldConfigOption = $timeseriesFieldConfigOption ?: "";
    }

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
