<?php

namespace Grafana\Foundation\Types\Dashboard;

class Panel implements \JsonSerializable {
    public string $title;

    public string $type;

    public ?\Grafana\Foundation\Types\Dashboard\DataSourceRef $datasource;

    /**
     * @var mixed
     */
    public $options;

    /**
     * @var array<\Grafana\Foundation\Runtime\Variants\Dataquery>
     */
    public ?array $targets;

    public ?\Grafana\Foundation\Types\Dashboard\FieldConfigSource $fieldConfig;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "title" => $this->title,
            "type" => $this->type,
        ];
        if (isset($this->datasource)) {
            $data["datasource"] = $this->datasource;
        }
        if (isset($this->options)) {
            $data["options"] = $this->options;
        }
        if (isset($this->targets)) {
            $data["targets"] = $this->targets;
        }
        if (isset($this->fieldConfig)) {
            $data["fieldConfig"] = $this->fieldConfig;
        }
        return $data;
    }
}
