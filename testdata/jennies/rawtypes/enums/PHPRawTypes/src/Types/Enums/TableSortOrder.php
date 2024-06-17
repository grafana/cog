<?php

namespace Grafana\Foundation\Types\Enums;

final class TableSortOrder implements \JsonSerializable, \Stringable {
    /**
     * @var string|int
     */
    private $value;

    /**
     * @var array<string, TableSortOrder>
     */
    private static $instances = [];

    private function __construct(string|int $value)
    {
        $this->value = $value;
    }

    public static function asc(): self
    {
        if (!isset(self::$instances["Asc"])) {
            self::$instances["Asc"] = new self("asc");
        }

        return self::$instances["Asc"];
    }

    public static function desc(): self
    {
        if (!isset(self::$instances["Desc"])) {
            self::$instances["Desc"] = new self("desc");
        }

        return self::$instances["Desc"];
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

