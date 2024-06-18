<?php

namespace Grafana\Foundation\VariantPanelcfgFull;

class Options implements \JsonSerializable {
    public string $timeseriesOption;

    /**
     * @param string|null $timeseriesOption
     */
    public function __construct(?string $timeseriesOption = null)
    {
        $this->timeseriesOption = $timeseriesOption ?: "";
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "timeseries_option" => $this->timeseriesOption,
        ];
        return $data;
    }
}
