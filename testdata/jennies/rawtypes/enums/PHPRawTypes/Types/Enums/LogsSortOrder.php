<?php

namespace Types\Enums;

final class LogsSortOrder implements \JsonSerializable, \Stringable {
    /**
     * @var string|int
     */
    private $value;

    /**
     * @var array<string, LogsSortOrder>
     */
    private static $instances = [];

    private function __construct(string|int $value)
    {
        $this->value = $value;
    }

    public function asc(): self
    {
        if (!isset(self::$instances["Asc"])) {
            self::$instances["Asc"] = new self("time_asc");
        }

        return self::$instances["Asc"];
    }

    public function desc(): self
    {
        if (!isset(self::$instances["Desc"])) {
            self::$instances["Desc"] = new self("time_desc");
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

