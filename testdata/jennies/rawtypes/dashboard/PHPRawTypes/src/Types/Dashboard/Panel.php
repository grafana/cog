<?php

namespace Grafana\Foundation\Types\Dashboard;

class Panel implements \JsonSerializable {
    public string $title;

    public string $type;

    public ?\Grafana\Foundation\Types\Dashboard\DataSourceRef $datasource;

    /**
     * @var mixed|null
     */
    public $options;

    /**
     * @var array<\Grafana\Foundation\Runtime\Variants\Dataquery>|null
     */
    public ?array $targets;

    public ?\Grafana\Foundation\Types\Dashboard\FieldConfigSource $fieldConfig;

    /**
     * @param string|null $title
     * @param string|null $type
     * @param \Grafana\Foundation\Types\Dashboard\DataSourceRef|null $datasource
     * @param mixed|null $options
     * @param array<\Grafana\Foundation\Runtime\Variants\Dataquery>|null $targets
     * @param \Grafana\Foundation\Types\Dashboard\FieldConfigSource|null $fieldConfig
     */
    public function __construct(?string $title = null, ?string $type = null, ?\Grafana\Foundation\Types\Dashboard\DataSourceRef $datasource = null,  $options = null, ?array $targets = null, ?\Grafana\Foundation\Types\Dashboard\FieldConfigSource $fieldConfig = null)
    {
        $this->title = $title ?: "";
        $this->type = $type ?: "";
        $this->datasource = $datasource;
        $this->options = $options;
        $this->targets = $targets;
        $this->fieldConfig = $fieldConfig;
    }

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
