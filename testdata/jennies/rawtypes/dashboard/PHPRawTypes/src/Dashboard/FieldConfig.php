<?php

namespace Grafana\Foundation\Dashboard;

class FieldConfig implements \JsonSerializable
{
    public ?string $unit;

    /**
     * @var mixed|null
     */
    public $custom;

    /**
     * @param string|null $unit
     * @param mixed|null $custom
     */
    public function __construct(?string $unit = null,  $custom = null)
    {
        $this->unit = $unit;
        $this->custom = $custom;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if (($this->unit === null) !== ($other->unit === null)) {
            return false;
        }
        if ($this->unit !== null) {
            if ($this->unit !== $other->unit) {
                return false;
            }
        }
    
        if (($this->custom === null) !== ($other->custom === null)) {
            return false;
        }
        if ($this->custom !== null) {
            if ($this->custom !== $other->custom) {
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
        /** @var array{unit?: string, custom?: mixed} $inputData */
        $data = $inputData;
        return new self(
            unit: $data["unit"] ?? null,
            custom: $data["custom"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->unit)) {
            $data->unit = $this->unit;
        }
        if (isset($this->custom)) {
            $data->custom = $this->custom;
        }
        return $data;
    }
}
