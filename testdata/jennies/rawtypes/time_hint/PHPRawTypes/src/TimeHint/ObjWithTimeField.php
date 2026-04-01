<?php

namespace Grafana\Foundation\TimeHint;

class ObjWithTimeField implements \JsonSerializable
{
    public string $registeredAt;

    public string $duration;

    /**
     * @param string|null $registeredAt
     * @param string|null $duration
     */
    public function __construct(?string $registeredAt = null, ?string $duration = null)
    {
        $this->registeredAt = $registeredAt ?: "";
        $this->duration = $duration ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{registeredAt?: string, duration?: string} $inputData */
        $data = $inputData;
        return new self(
            registeredAt: $data["registeredAt"] ?? null,
            duration: $data["duration"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->registeredAt = $this->registeredAt;
        $data->duration = $this->duration;
        return $data;
    }
}
