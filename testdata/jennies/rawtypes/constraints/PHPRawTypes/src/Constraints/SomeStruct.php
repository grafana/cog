<?php

namespace Grafana\Foundation\Constraints;

class SomeStruct implements \JsonSerializable
{
    public int $id;

    public ?int $maybeId;

    public int $greaterThanZero;

    public int $negative;

    public string $title;

    /**
     * @var array<string, string>
     */
    public array $labels;

    /**
     * @var array<string>
     */
    public array $tags;

    public string $regex;

    public string $negativeRegex;

    /**
     * @param int|null $id
     * @param int|null $maybeId
     * @param int|null $greaterThanZero
     * @param int|null $negative
     * @param string|null $title
     * @param array<string, string>|null $labels
     * @param array<string>|null $tags
     * @param string|null $regex
     * @param string|null $negativeRegex
     */
    public function __construct(?int $id = null, ?int $maybeId = null, ?int $greaterThanZero = null, ?int $negative = null, ?string $title = null, ?array $labels = null, ?array $tags = null, ?string $regex = null, ?string $negativeRegex = null)
    {
        $this->id = $id ?: 0;
        $this->maybeId = $maybeId;
        $this->greaterThanZero = $greaterThanZero ?: 0;
        $this->negative = $negative ?: 0;
        $this->title = $title ?: "";
        $this->labels = $labels ?: [];
        $this->tags = $tags ?: [];
        $this->regex = $regex ?: "";
        $this->negativeRegex = $negativeRegex ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{id?: int, maybeId?: int, greaterThanZero?: int, negative?: int, title?: string, labels?: array<string, string>, tags?: array<string>, regex?: string, negativeRegex?: string} $inputData */
        $data = $inputData;
        return new self(
            id: $data["id"] ?? null,
            maybeId: $data["maybeId"] ?? null,
            greaterThanZero: $data["greaterThanZero"] ?? null,
            negative: $data["negative"] ?? null,
            title: $data["title"] ?? null,
            labels: $data["labels"] ?? null,
            tags: $data["tags"] ?? null,
            regex: $data["regex"] ?? null,
            negativeRegex: $data["negativeRegex"] ?? null,
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
        $data->negative = $this->negative;
        $data->title = $this->title;
        $data->labels = $this->labels;
        $data->tags = $this->tags;
        $data->regex = $this->regex;
        $data->negativeRegex = $this->negativeRegex;
        if (isset($this->maybeId)) {
            $data->maybeId = $this->maybeId;
        }
        return $data;
    }
}
