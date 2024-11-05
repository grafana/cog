<?php

namespace Grafana\Foundation\TimeHint;

class ObjWithTimeField implements \JsonSerializable
{
    public string $registeredAt;

    /**
     * @param string|null $registeredAt
     */
    public function __construct(?string $registeredAt = null)
    {
        $this->registeredAt = $registeredAt ?: "";
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{registeredAt?: string} $inputData */
        $data = $inputData;
        return new self(
            registeredAt: $data["registeredAt"] ?? null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "registeredAt" => $this->registeredAt,
        ];
        return $data;
    }
}
