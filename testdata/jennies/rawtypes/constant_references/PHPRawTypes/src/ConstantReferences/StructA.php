<?php

namespace Grafana\Foundation\ConstantReferences;

class StructA implements \JsonSerializable
{
    public \Grafana\Foundation\ConstantReferences\Enum $myEnum;

    public ?\Grafana\Foundation\ConstantReferences\Enum $other;

    public function __construct()
    {
        $this->myEnum = \Grafana\Foundation\ConstantReferences\Enum::valueA();
        $this->other = \Grafana\Foundation\ConstantReferences\Enum::valueA();
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if ($this->myEnum !== $other->myEnum) {
            return false;
        }
    
        if (($this->other === null) !== ($other->other === null)) {
            return false;
        }
        if ($this->other !== null) {
            if ($this->other !== $other->other) {
                return false;
            }
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        return new self(
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->myEnum = $this->myEnum;
        if (isset($this->other)) {
            $data->other = $this->other;
        }
        return $data;
    }
}
