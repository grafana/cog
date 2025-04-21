<?php

namespace Grafana\Foundation\Constraints;

class RefStruct implements \JsonSerializable
{
    /**
     * @var array<string, string>
     */
    public array $labels;

    /**
     * @var array<string>
     */
    public array $tags;

    /**
     * @param array<string, string>|null $labels
     * @param array<string>|null $tags
     */
    public function __construct(?array $labels = null, ?array $tags = null)
    {
        $this->labels = $labels ?: [];
        $this->tags = $tags ?: [];
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{labels?: array<string, string>, tags?: array<string>} $inputData */
        $data = $inputData;
        return new self(
            labels: $data["labels"] ?? null,
            tags: $data["tags"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->labels = $this->labels;
        $data->tags = $this->tags;
        return $data;
    }
}
