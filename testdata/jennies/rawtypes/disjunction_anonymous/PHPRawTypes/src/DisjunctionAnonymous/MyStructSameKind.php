<?php

namespace Grafana\Foundation\DisjunctionAnonymous;

final class MyStructSameKind implements \JsonSerializable, \Stringable {
    /**
     * @var string
     */
    private $value;

    /**
     * @var array<string, MyStructSameKind>
     */
    private static $instances = [];

    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function a(): self
    {
        if (!isset(self::$instances["A"])) {
            self::$instances["A"] = new self("a");
        }

        return self::$instances["A"];
    }

    public static function b(): self
    {
        if (!isset(self::$instances["B"])) {
            self::$instances["B"] = new self("b");
        }

        return self::$instances["B"];
    }

    public static function c(): self
    {
        if (!isset(self::$instances["C"])) {
            self::$instances["C"] = new self("c");
        }

        return self::$instances["C"];
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

        throw new \UnexpectedValueException("Value '$value' is not part of the enum MyStructSameKind");
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

