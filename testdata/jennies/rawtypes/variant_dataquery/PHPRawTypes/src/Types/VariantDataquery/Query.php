<?php

namespace Grafana\Foundation\Types\VariantDataquery;

class Query implements \Grafana\Foundation\Runtime\Variants\Dataquery
{
    public string $expr;

    public ?bool $instant;
}
