<?php

namespace Types\VariantDataquery;

class Query implements \Runtime\Variants\Dataquery
{
    public string $expr;

    public ?bool $instant;
}
