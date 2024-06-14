<?php

namespace Types\Enums;

/**
 * This is a very interesting string enum.
 */
final class Operator implements \JsonSerializable, \Stringable {
    /**
     * @var string|int
     */
    private $value;

    /**
     * @var array<string, Operator>
     */
    private static $instances = [];

    private function __construct(string|int $value)
    {
        $this->value = $value;
    }

    public function greaterThan(): self
    {
        if (!isset(self::$instances["GreaterThan"])) {
            self::$instances["GreaterThan"] = new self(">");
        }

        return self::$instances["GreaterThan"];
    }

    public function lessThan(): self
    {
        if (!isset(self::$instances["LessThan"])) {
            self::$instances["LessThan"] = new self("<");
        }

        return self::$instances["LessThan"];
    }

    public function jsonSerialize(): string|int
    {
        return $this->value;
    }

    public function __toString(): string
    {
        return (string) $this->value;
    }
}

