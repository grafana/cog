<?php

namespace Grafana\Foundation\Types\Dashboard;

class Dashboard implements \JsonSerializable {
    public string $title;

    /**
     * @var array<\Grafana\Foundation\Types\Dashboard\Panel>
     */
    public ?array $panels;

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
