<?php

namespace Grafana\Foundation\Dashboard;

class Panel implements \JsonSerializable
{
    public string $title;

    public string $type;

    public ?\Grafana\Foundation\Dashboard\DataSourceRef $datasource;

    /**
     * @var mixed|null
     */
    public $options;

    /**
     * @var array<\Grafana\Foundation\Cog\Dataquery>|null
     */
    public ?array $targets;

    public ?\Grafana\Foundation\Dashboard\FieldConfigSource $fieldConfig;

    /**
     * @param string|null $title
     * @param string|null $type
     * @param \Grafana\Foundation\Dashboard\DataSourceRef|null $datasource
     * @param mixed|null $options
     * @param array<\Grafana\Foundation\Cog\Dataquery>|null $targets
     * @param \Grafana\Foundation\Dashboard\FieldConfigSource|null $fieldConfig
     */
    public function __construct(?string $title = null, ?string $type = null, ?\Grafana\Foundation\Dashboard\DataSourceRef $datasource = null,  $options = null, ?array $targets = null, ?\Grafana\Foundation\Dashboard\FieldConfigSource $fieldConfig = null)
    {
        $this->title = $title ?: "";
        $this->type = $type ?: "";
        $this->datasource = $datasource;
        $this->options = $options;
        $this->targets = $targets;
        $this->fieldConfig = $fieldConfig;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{title?: string, type?: string, datasource?: mixed, options?: mixed, targets?: array<mixed>, fieldConfig?: mixed} $inputData */
        $data = $inputData;
        return new self(
            title: $data["title"] ?? null,
            type: $data["type"] ?? null,
            datasource: isset($data["datasource"]) ? (function($input) {
    	/** @var array{type?: string, uid?: string} */
    $val = $input;
    	return \Grafana\Foundation\Dashboard\DataSourceRef::fromArray($val);
    })($data["datasource"]) : null,
            options: $data["options"] ?? null,
            targets: can not generate custom unmarshal function for composable slot with variant 'dataquery': template block variant_dataquery_field_unmarshal not found,
            fieldConfig: isset($data["fieldConfig"]) ? (function($input) {
    	/** @var array{defaults?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\Dashboard\FieldConfigSource::fromArray($val);
    })($data["fieldConfig"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->title = $this->title;
        $data->type = $this->type;
        if (isset($this->datasource)) {
            $data->datasource = $this->datasource;
        }
        if (isset($this->options)) {
            $data->options = $this->options;
        }
        if (isset($this->targets)) {
            $data->targets = $this->targets;
        }
        if (isset($this->fieldConfig)) {
            $data->fieldConfig = $this->fieldConfig;
        }
        return $data;
    }
}
