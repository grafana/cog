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

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->uid !== $other->uid) {
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
        /** @var array{uid?: string, intVal?: int} $inputData */
        $data = $inputData;
        return new self(
            uid: $data["uid"] ?? null,
            intVal: $data["intVal"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->uid = $this->uid;
        $data->intVal = $this->intVal;
        return $data;
    }
}
