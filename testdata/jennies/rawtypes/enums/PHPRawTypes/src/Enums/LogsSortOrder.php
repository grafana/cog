<?php

namespace Grafana\Foundation\Enums;

final class LogsSortOrder implements \JsonSerializable, \Stringable {
    /**
     * @var string
     */
    private $value;

    /**
     * @var array<string, LogsSortOrder>
     */
    private static $instances = [];

    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function asc(): self
    {
        if (!isset(self::$instances["Asc"])) {
            self::$instances["Asc"] = new self("time_asc");
        }

        return self::$instances["Asc"];
    }

    public static function desc(): self
    {
        if (!isset(self::$instances["Desc"])) {
            self::$instances["Desc"] = new self("time_desc");
        }

        return self::$instances["Desc"];
    }

    public static function fromValue(string $value): self
    {
        if ($value === "time_asc") {
            return self::asc();
        }

        if ($value === "time_desc") {
            return self::desc();
        }

        throw new \UnexpectedValueException("Value '$value' is not part of the enum LogsSortOrder");
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

