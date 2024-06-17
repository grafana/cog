<?php

namespace Grafana\Foundation\Types\VariantPanelcfgFull;

class Options implements \JsonSerializable {
    public string $timeseriesOption;

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
