<?php

namespace Grafana\Foundation\Enums;

/**
 * This is a very interesting string enum.
 */
final class Operator implements \JsonSerializable, \Stringable {
    /**
     * @var string
     */
    private $value;

    /**
     * @var array<string, Operator>
     */
    private static $instances = [];

    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function greaterThan(): self
    {
        if (!isset(self::$instances["GreaterThan"])) {
            self::$instances["GreaterThan"] = new self(">");
        }

        return self::$instances["GreaterThan"];
    }

    public static function lessThan(): self
    {
        if (!isset(self::$instances["LessThan"])) {
            self::$instances["LessThan"] = new self("<");
        }

        return self::$instances["LessThan"];
    }

    public static function fromValue(string $value): self
    {
        if ($value === ">") {
            return self::greaterThan();
        }

        if ($value === "<") {
            return self::lessThan();
        }

        throw new \UnexpectedValueException("Value '$value' is not part of the enum Operator");
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

