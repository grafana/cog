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

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->nestedVal !== $other->nestedVal) {
            return false;
        }
    
        return true;
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
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->nestedVal = $this->nestedVal;
        return $data;
    }
}
