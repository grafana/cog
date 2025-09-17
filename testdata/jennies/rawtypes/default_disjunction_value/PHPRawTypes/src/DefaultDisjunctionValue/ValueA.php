<?php

namespace Grafana\Foundation\DefaultDisjunctionValue;

class ValueA implements \JsonSerializable
{
    public string $type;

    /**
     * @var array<string>
     */
    public array $anArray;

    public \Grafana\Foundation\DefaultDisjunctionValue\ValueB $otherRef;

    /**
     * @param array<string>|null $anArray
     * @param \Grafana\Foundation\DefaultDisjunctionValue\ValueB|null $otherRef
     */
    public function __construct(?array $anArray = null, ?\Grafana\Foundation\DefaultDisjunctionValue\ValueB $otherRef = null)
    {
        $this->type = "A";
    
        $this->anArray = $anArray ?: [];
        $this->otherRef = $otherRef ?: new \Grafana\Foundation\DefaultDisjunctionValue\ValueB();
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{type?: string, anArray?: array<string>, otherRef?: mixed} $inputData */
        $data = $inputData;
        return new self(
            anArray: $data["anArray"] ?? null,
            otherRef: isset($data["otherRef"]) ? (function($input) {
    	/** @var array{type?: string, aMap?: array<string, int>, def?: int|string|bool} */
    $val = $input;
    	return \Grafana\Foundation\DefaultDisjunctionValue\ValueB::fromArray($val);
    })($data["otherRef"]) : null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
            "type" => $this->type,
            "anArray" => $this->anArray,
            "otherRef" => $this->otherRef,
        ];
        return $data;
    }
}
