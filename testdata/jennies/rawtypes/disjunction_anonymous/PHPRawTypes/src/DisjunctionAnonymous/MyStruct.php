<?php

namespace Grafana\Foundation\DisjunctionAnonymous;

class MyStruct implements \JsonSerializable
{
    /**
     * @var string|bool|float|int
     */
    public $scalars;

    public \Grafana\Foundation\DisjunctionAnonymous\MyStructSameKind $sameKind;

    /**
     * @var mixed
     */
    public $refs;

    /**
     * @var \Grafana\Foundation\DisjunctionAnonymous\StructA|string|int
     */
    public $mixed;

    /**
     * @param string|bool|float|int|null $scalars
     * @param \Grafana\Foundation\DisjunctionAnonymous\MyStructSameKind|null $sameKind
     * @param mixed|null $refs
     * @param \Grafana\Foundation\DisjunctionAnonymous\StructA|string|int|null $mixed
     */
    public function __construct( $scalars = null, ?\Grafana\Foundation\DisjunctionAnonymous\MyStructSameKind $sameKind = null,  $refs = null,  $mixed = null)
    {
        $this->scalars = $scalars ?: "";
        $this->sameKind = $sameKind ?: \Grafana\Foundation\DisjunctionAnonymous\MyStructSameKind::A();
        $this->refs = $refs ?: null;
        $this->mixed = $mixed ?: new \Grafana\Foundation\DisjunctionAnonymous\StructA();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{scalars?: string|bool|float|int, sameKind?: string, refs?: mixed, mixed?: mixed|string|int} $inputData */
        $data = $inputData;
        return new self(
            scalars: isset($data["scalars"]) ? (function($input) {
        switch (true) {
        case is_string($input):
            return $input;
        case is_bool($input):
            return $input;
        case is_float($input):
            return $input;
        case is_int($input):
            return $input;
        default:
            throw new \ValueError('incorrect value for disjunction');
    }
    })($data["scalars"]) : null,
            sameKind: isset($data["sameKind"]) ? (function($input) { return \Grafana\Foundation\DisjunctionAnonymous\MyStructSameKind::fromValue($input); })($data["sameKind"]) : null,
            refs: $data["refs"] ?? null,
            mixed: isset($data["mixed"]) ? (function($input) {
        switch (true) {
        case is_string($input):
            return $input;
        case is_int($input):
            return $input;
        default:
            /** @var array{field?: string} $input */
            return \Grafana\Foundation\DisjunctionAnonymous\StructA::fromArray($input);
    }
    })($data["mixed"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->scalars = $this->scalars;
        $data->sameKind = $this->sameKind;
        $data->refs = $this->refs;
        $data->mixed = $this->mixed;
        return $data;
    }
}
