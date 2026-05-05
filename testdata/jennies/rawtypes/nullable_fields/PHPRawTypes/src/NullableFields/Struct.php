<?php

namespace Grafana\Foundation\NullableFields;

class Struct implements \JsonSerializable
{
    public ?\Grafana\Foundation\NullableFields\MyObject $a;

    public ?\Grafana\Foundation\NullableFields\MyObject $b;

    public ?string $c;

    /**
     * @var array<string>|null
     */
    public ?array $d;

    /**
     * @var array<string, string|null>
     */
    public array $e;

    public ?\Grafana\Foundation\NullableFields\NullableFieldsStructF $f;

    public ?string $g;

    /**
     * @param \Grafana\Foundation\NullableFields\MyObject|null $a
     * @param \Grafana\Foundation\NullableFields\MyObject|null $b
     * @param string|null $c
     * @param array<string>|null $d
     * @param array<string, string|null>|null $e
     * @param \Grafana\Foundation\NullableFields\NullableFieldsStructF|null $f
     */
    public function __construct(?\Grafana\Foundation\NullableFields\MyObject $a = null, ?\Grafana\Foundation\NullableFields\MyObject $b = null, ?string $c = null, ?array $d = null, ?array $e = null, ?\Grafana\Foundation\NullableFields\NullableFieldsStructF $f = null)
    {
        $this->a = $a;
        $this->b = $b;
        $this->c = $c;
        $this->d = $d;
        $this->e = $e ?: [];
        $this->f = $f;
        $this->g = \Grafana\Foundation\NullableFields\Constants::CONSTANT_REF;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if (($this->a === null) !== ($other->a === null)) {
            return false;
        }
        if ($this->a !== null) {
            if (!$this->a->equals($other->a)) {
                return false;
            }
        }
    
        if (($this->b === null) !== ($other->b === null)) {
            return false;
        }
        if ($this->b !== null) {
            if (!$this->b->equals($other->b)) {
                return false;
            }
        }
    
        if (($this->c === null) !== ($other->c === null)) {
            return false;
        }
        if ($this->c !== null) {
            if ($this->c !== $other->c) {
                return false;
            }
        }
    
        if (($this->d === null) !== ($other->d === null)) {
            return false;
        }
        if ($this->d !== null) {
            if ($this->d != $other->d) {
                return false;
            }
        }
    
        if ($this->e != $other->e) {
            return false;
        }
    
        if (($this->f === null) !== ($other->f === null)) {
            return false;
        }
        if ($this->f !== null) {
            if (!$this->f->equals($other->f)) {
                return false;
            }
        }
    
        if (($this->g === null) !== ($other->g === null)) {
            return false;
        }
        if ($this->g !== null) {
            if ($this->g !== $other->g) {
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
        /** @var array{a?: mixed, b?: mixed, c?: string, d?: array<string>, e?: array<string, string>, f?: mixed, g?: "hey"} $inputData */
        $data = $inputData;
        return new self(
            a: isset($data["a"]) ? (function($input) {
    	/** @var array{field?: string} */
    $val = $input;
    	return \Grafana\Foundation\NullableFields\MyObject::fromArray($val);
    })($data["a"]) : null,
            b: isset($data["b"]) ? (function($input) {
    	/** @var array{field?: string} */
    $val = $input;
    	return \Grafana\Foundation\NullableFields\MyObject::fromArray($val);
    })($data["b"]) : null,
            c: $data["c"] ?? null,
            d: $data["d"] ?? null,
            e: $data["e"] ?? null,
            f: isset($data["f"]) ? (function($input) {
    	/** @var array{a?: string} */
    $val = $input;
    	return \Grafana\Foundation\NullableFields\NullableFieldsStructF::fromArray($val);
    })($data["f"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->e = $this->e;
        if (isset($this->a)) {
            $data->a = $this->a;
        }
        if (isset($this->b)) {
            $data->b = $this->b;
        }
        if (isset($this->c)) {
            $data->c = $this->c;
        }
        if (isset($this->d)) {
            $data->d = $this->d;
        }
        if (isset($this->f)) {
            $data->f = $this->f;
        }
        if (isset($this->g)) {
            $data->g = $this->g;
        }
        return $data;
    }
}
