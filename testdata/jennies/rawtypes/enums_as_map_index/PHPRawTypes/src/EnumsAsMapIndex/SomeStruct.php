<?php

namespace Grafana\Foundation\EnumsAsMapIndex;

class SomeStruct implements \JsonSerializable
{
    /**
     * @var array<\Grafana\Foundation\EnumsAsMapIndex\StringEnum, string>
     */
    public array $data;

    /**
     * @param array<\Grafana\Foundation\EnumsAsMapIndex\StringEnum, string>|null $data
     */
    public function __construct(?array $data = null)
    {
        $this->data = $data ?: [];
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->data != $other->data) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{data?: array<string, string>} $inputData */
        $data = $inputData;
        return new self(
            data: $data["data"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->data = $this->data;
        return $data;
    }
}
