<?php

namespace Grafana\Foundation\VariantPanelcfgFull;

class Options implements \JsonSerializable
{
    public string $timeseriesOption;

    /**
     * @param string|null $timeseriesOption
     */
    public function __construct(?string $timeseriesOption = null)
    {
        $this->timeseriesOption = $timeseriesOption ?: "";
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->timeseriesOption !== $other->timeseriesOption) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{timeseries_option?: string} $inputData */
        $data = $inputData;
        return new self(
            timeseriesOption: $data["timeseries_option"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->timeseries_option = $this->timeseriesOption;
        return $data;
    }
}
