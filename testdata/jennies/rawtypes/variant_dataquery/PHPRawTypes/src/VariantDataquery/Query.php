<?php

namespace Grafana\Foundation\VariantDataquery;

class Query implements \JsonSerializable, \Grafana\Foundation\Runtime\Dataquery {
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
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "expr" => $this->expr,
        ];
        if (isset($this->instant)) {
            $data["instant"] = $this->instant;
        }
        return $data;
    }
}
