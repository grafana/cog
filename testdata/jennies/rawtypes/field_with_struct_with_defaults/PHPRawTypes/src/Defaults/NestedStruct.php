<?php

namespace Grafana\Foundation\Defaults;

class NestedStruct implements \JsonSerializable
{
    public string $stringVal;

    public int $intVal;

    /**
     * @param string|null $stringVal
     * @param int|null $intVal
     */
    public function __construct(?string $stringVal = null, ?int $intVal = null)
    {
        $this->stringVal = $stringVal ?: "";
        $this->intVal = $intVal ?: 0;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->stringVal !== $other->stringVal) {
            return false;
        }
    
        if ($this->intVal !== $other->intVal) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{stringVal?: string, intVal?: int} $inputData */
        $data = $inputData;
        return new self(
            stringVal: $data["stringVal"] ?? null,
            intVal: $data["intVal"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->stringVal = $this->stringVal;
        $data->intVal = $this->intVal;
        return $data;
    }
}
