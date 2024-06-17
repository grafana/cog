<?php

namespace Grafana\Foundation\Types\VariantPanelcfgOnlyOptions;

class Options implements \JsonSerializable {
    public string $content;

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "content" => $this->content,
        ];
        return $data;
    }
}
