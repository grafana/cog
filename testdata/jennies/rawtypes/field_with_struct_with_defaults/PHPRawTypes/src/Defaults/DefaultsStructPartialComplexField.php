<?php

namespace Grafana\Foundation\Defaults;

class DefaultsStructPartialComplexField implements \JsonSerializable
{
    public string $uid;

    public int $intVal;

    /**
     * @param string|null $uid
     * @param int|null $intVal
     */
    public function __construct(?string $uid = null, ?int $intVal = null)
    {
        $this->uid = $uid ?: "";
        $this->intVal = $intVal ?: 0;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{uid?: string, intVal?: int} $inputData */
        $data = $inputData;
        return new self(
            uid: $data["uid"] ?? null,
            intVal: $data["intVal"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "uid" => $this->uid,
            "intVal" => $this->intVal,
        ];
        return $data;
    }
}
