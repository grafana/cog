<?php

namespace Grafana\Foundation\VariantDataquery;

class Query implements \JsonSerializable, \Grafana\Foundation\Cog\Dataquery
{
    public string $expr;

    public ?bool $instant;

    /**
     * @param string|null $expr
     * @param bool|null $instant
     */
    public function __construct(?string $expr = null, ?bool $instant = null)
    {
        $this->expr = $expr ?: "";
        $this->instant = $instant;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{expr?: string, instant?: bool} $inputData */
        $data = $inputData;
        return new self(
            expr: $data["expr"] ?? null,
            instant: $data["instant"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->expr = $this->expr;
        if (isset($this->instant)) {
            $data->instant = $this->instant;
        }
        return $data;
    }
}
