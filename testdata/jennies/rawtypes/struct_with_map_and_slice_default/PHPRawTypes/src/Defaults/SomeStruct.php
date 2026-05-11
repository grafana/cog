<?php

namespace Grafana\Foundation\Defaults;

class SomeStruct implements \JsonSerializable
{
    /**
     * @var array<string, mixed>|null
     */
    public ?array $options;

    /**
     * @var array<string>|null
     */
    public ?array $items;

    /**
     * @var mixed
     */
    public $extra;

    /**
     * @param array<string, mixed>|null $options
     * @param array<string>|null $items
     * @param mixed|null $extra
     */
    public function __construct(?array $options = null, ?array $items = null,  $extra = null)
    {
        $this->options = $options;
        $this->items = $items;
        $this->extra = $extra ?: [];
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if (($this->options === null) !== ($other->options === null)) {
            return false;
        }
        if ($this->options !== null) {
            if ($this->options != $other->options) {
                return false;
            }
        }
    
        if (($this->items === null) !== ($other->items === null)) {
            return false;
        }
        if ($this->items !== null) {
            if ($this->items != $other->items) {
                return false;
            }
        }
    
        if ($this->extra !== $other->extra) {
            return false;
        }
    
        return true;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{options?: array<string, mixed>, items?: array<string>, extra?: mixed} $inputData */
        $data = $inputData;
        return new self(
            options: $data["options"] ?? null,
            items: $data["items"] ?? null,
            extra: $data["extra"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->extra = $this->extra;
        if (isset($this->options)) {
            $data->options = $this->options;
        }
        if (isset($this->items)) {
            $data->items = $this->items;
        }
        return $data;
    }
}
