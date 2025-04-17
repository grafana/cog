<?php

namespace Grafana\Foundation\Constraints;

class SomeStruct implements \JsonSerializable
{
    public int $id;

    public ?int $maybeId;

    public string $title;

    public ?\Grafana\Foundation\Constraints\RefStruct $refStruct;

    /**
     * @param int|null $id
     * @param int|null $maybeId
     * @param string|null $title
     * @param \Grafana\Foundation\Constraints\RefStruct|null $refStruct
     */
    public function __construct(?int $id = null, ?int $maybeId = null, ?string $title = null, ?\Grafana\Foundation\Constraints\RefStruct $refStruct = null)
    {
        $this->id = $id ?: 0;
        $this->maybeId = $maybeId;
        $this->title = $title ?: "";
        $this->refStruct = $refStruct;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{id?: int, maybeId?: int, title?: string, refStruct?: mixed} $inputData */
        $data = $inputData;
        return new self(
            id: $data["id"] ?? null,
            maybeId: $data["maybeId"] ?? null,
            title: $data["title"] ?? null,
            refStruct: isset($data["refStruct"]) ? (function($input) {
    	/** @var array{labels?: array<string, string>, tags?: array<string>} */
    $val = $input;
    	return \Grafana\Foundation\Constraints\RefStruct::fromArray($val);
    })($data["refStruct"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->id = $this->id;
        $data->title = $this->title;
        if (isset($this->maybeId)) {
            $data->maybeId = $this->maybeId;
        }
        if (isset($this->refStruct)) {
            $data->refStruct = $this->refStruct;
        }
        return $data;
    }
}
