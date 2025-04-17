<?php

namespace Grafana\Foundation\VariantPanelcfgFull;

class FieldConfig implements \JsonSerializable
{
    public string $timeseriesFieldConfigOption;

    /**
     * @param string|null $timeseriesFieldConfigOption
     */
    public function __construct(?string $timeseriesFieldConfigOption = null)
    {
        $this->timeseriesFieldConfigOption = $timeseriesFieldConfigOption ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{timeseries_field_config_option?: string} $inputData */
        $data = $inputData;
        return new self(
            timeseriesFieldConfigOption: $data["timeseries_field_config_option"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->timeseries_field_config_option = $this->timeseriesFieldConfigOption;
        return $data;
    }
}
