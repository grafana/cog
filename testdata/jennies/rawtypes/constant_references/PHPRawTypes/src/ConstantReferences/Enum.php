<?php

namespace Grafana\Foundation\ConstantReferences;

final class Enum implements \JsonSerializable, \Stringable {
    /**
     * @var string
     */
    private $value;

    /**
     * @var array<string, Enum>
     */
    private static $instances = [];

    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function valueA(): self
    {
        if (!isset(self::$instances["ValueA"])) {
            self::$instances["ValueA"] = new self("ValueA");
        }

        return self::$instances["ValueA"];
    }

    public static function valueB(): self
    {
        if (!isset(self::$instances["ValueB"])) {
            self::$instances["ValueB"] = new self("ValueB");
        }

        return self::$instances["ValueB"];
    }

    public static function valueC(): self
    {
        if (!isset(self::$instances["ValueC"])) {
            self::$instances["ValueC"] = new self("ValueC");
        }

        return self::$instances["ValueC"];
    }

    public static function fromValue(string $value): self
    {
        if ($value === "ValueA") {
            return self::valueA();
        }

        if ($value === "ValueB") {
            return self::valueB();
        }

        if ($value === "ValueC") {
            return self::valueC();
        }

        throw new \UnexpectedValueException("Value '$value' is not part of the enum Enum");
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

