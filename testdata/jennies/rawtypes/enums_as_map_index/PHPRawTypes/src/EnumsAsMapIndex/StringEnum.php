<?php

namespace Grafana\Foundation\EnumsAsMapIndex;

final class StringEnum implements \JsonSerializable, \Stringable {
    /**
     * @var string
     */
    private $value;

    /**
     * @var array<string, StringEnum>
     */
    private static $instances = [];

    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function a(): self
    {
        if (!isset(self::$instances["a"])) {
            self::$instances["a"] = new self("a");
        }

        return self::$instances["a"];
    }

    public static function b(): self
    {
        if (!isset(self::$instances["b"])) {
            self::$instances["b"] = new self("b");
        }

        return self::$instances["b"];
    }

    public static function c(): self
    {
        if (!isset(self::$instances["c"])) {
            self::$instances["c"] = new self("c");
        }

        return self::$instances["c"];
    }

    public static function fromValue(string $value): self
    {
        if ($value === "a") {
            return self::a();
        }

        if ($value === "b") {
            return self::b();
        }

        if ($value === "c") {
            return self::c();
        }

        throw new \UnexpectedValueException("Value '$value' is not part of the enum StringEnum");
    }

    public function jsonSerialize(): string
    {
        return $this->value;
    }

    public function __toString(): string
    {
        return $this->value;
    }
}

