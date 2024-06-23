<?php

namespace Grafana\Foundation\Defaults;

class DefaultsStructComplexFieldNested implements \JsonSerializable
{
    public string $nestedVal;

    /**
     * @param string|null $nestedVal
     */
    public function __construct(?string $nestedVal = null)
    {
        $this->nestedVal = $nestedVal ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{nestedVal?: string} $inputData */
        $data = $inputData;
        return new self(
            nestedVal: $data["nestedVal"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "nestedVal" => $this->nestedVal,
        ];
        return $data;
    }
}
