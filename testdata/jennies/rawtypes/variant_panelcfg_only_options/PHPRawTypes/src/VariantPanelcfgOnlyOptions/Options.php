<?php

namespace Grafana\Foundation\VariantPanelcfgOnlyOptions;

class Options implements \JsonSerializable {
    public string $content;

    /**
     * @param string|null $content
     */
    public function __construct(?string $content = null)
    {
        $this->content = $content ?: "";
    }

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
