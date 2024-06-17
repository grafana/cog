<?php

namespace Grafana\Foundation\Types\VariantDataquery;

class Query implements \JsonSerializable, \Grafana\Foundation\Runtime\Variants\Dataquery {
    public string $expr;

    public ?bool $instant;

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
