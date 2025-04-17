<?php

namespace Grafana\Foundation\Dashboard;

class Dashboard implements \JsonSerializable
{
    public string $title;

    /**
     * @var array<\Grafana\Foundation\Dashboard\Panel>|null
     */
    public ?array $panels;

    /**
     * @param string|null $title
     * @param array<\Grafana\Foundation\Dashboard\Panel>|null $panels
     */
    public function __construct(?string $title = null, ?array $panels = null)
    {
        $this->title = $title ?: "";
        $this->panels = $panels;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{title?: string, panels?: array<mixed>} $inputData */
        $data = $inputData;
        return new self(
            title: $data["title"] ?? null,
            panels: array_filter(array_map((function($input) {
    	/** @var array{title?: string, type?: string, datasource?: mixed, options?: mixed, targets?: array<mixed>, fieldConfig?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\Dashboard\Panel::fromArray($val);
    }), $data["panels"] ?? [])),
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->title = $this->title;
        if (isset($this->panels)) {
            $data->panels = $this->panels;
        }
        return $data;
    }
}
