<?php

namespace Grafana\Foundation\DefaultDisjunctionValue;

class ValueB implements \JsonSerializable
{
    public string $type;

    /**
     * @var array<string, int>
     */
    public array $aMap;

    /**
     * @var int|string|bool
     */
    public $def;

    /**
     * @param array<string, int>|null $aMap
     * @param int|string|bool|null $def
     */
    public function __construct(?array $aMap = null,  $def = null)
    {
        $this->type = "B";
    
        $this->aMap = $aMap ?: [];
        $this->def = $def ?: 1;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{type?: string, aMap?: array<string, int>, def?: int|string|bool} $inputData */
        $data = $inputData;
        return new self(
            aMap: $data["aMap"] ?? null,
            def: isset($data["def"]) ? (function($input) {
        switch (true) {
        case is_int($input):
            return $input;
        case is_string($input):
            return $input;
        case is_bool($input):
            return $input;
        default:
            throw new \ValueError('incorrect value for disjunction');
    }
    })($data["def"]) : null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "type" => $this->type,
            "aMap" => $this->aMap,
            "def" => $this->def,
        ];
        return $data;
    }
}
