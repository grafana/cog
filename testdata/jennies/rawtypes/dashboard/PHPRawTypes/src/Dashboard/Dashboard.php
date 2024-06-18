<?php

namespace Grafana\Foundation\Dashboard;

class Dashboard implements \JsonSerializable {
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
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "title" => $this->title,
        ];
        if (isset($this->panels)) {
            $data["panels"] = $this->panels;
        }
        return $data;
    }
}
