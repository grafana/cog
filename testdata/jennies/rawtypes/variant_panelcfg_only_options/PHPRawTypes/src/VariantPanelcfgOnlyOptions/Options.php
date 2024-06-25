<?php

namespace Grafana\Foundation\VariantPanelcfgOnlyOptions;

class Options implements \JsonSerializable
{
    public string $content;

    /**
     * @param string|null $content
     */
    public function __construct(?string $content = null)
    {
        $this->content = $content ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{content?: string} $inputData */
        $data = $inputData;
        return new self(
            content: $data["content"] ?? null,
        );
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
