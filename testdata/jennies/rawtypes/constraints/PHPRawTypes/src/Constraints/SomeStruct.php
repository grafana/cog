<?php

namespace Grafana\Foundation\Constraints;

class SomeStruct implements \JsonSerializable
{
    public int $id;

    public ?int $maybeId;

    public int $greaterThanZero;

    public string $title;

    /**
     * @var array<string, string>
     */
    public array $labels;

    /**
     * @var array<string>
     */
    public array $tags;

    /**
     * @param int|null $id
     * @param int|null $maybeId
     * @param int|null $greaterThanZero
     * @param string|null $title
     * @param array<string, string>|null $labels
     * @param array<string>|null $tags
     */
    public function __construct(?int $id = null, ?int $maybeId = null, ?int $greaterThanZero = null, ?string $title = null, ?array $labels = null, ?array $tags = null)
    {
        $this->id = $id ?: 0;
        $this->maybeId = $maybeId;
        $this->greaterThanZero = $greaterThanZero ?: 0;
        $this->title = $title ?: "";
        $this->labels = $labels ?: [];
        $this->tags = $tags ?: [];
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{id?: int, maybeId?: int, greaterThanZero?: int, title?: string, labels?: array<string, string>, tags?: array<string>} $inputData */
        $data = $inputData;
        return new self(
            id: $data["id"] ?? null,
            maybeId: $data["maybeId"] ?? null,
            greaterThanZero: $data["greaterThanZero"] ?? null,
            title: $data["title"] ?? null,
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
        $data->id = $this->id;
        $data->greaterThanZero = $this->greaterThanZero;
        $data->title = $this->title;
        $data->labels = $this->labels;
        $data->tags = $this->tags;
        if (isset($this->maybeId)) {
            $data->maybeId = $this->maybeId;
        }
        return $data;
    }
}
