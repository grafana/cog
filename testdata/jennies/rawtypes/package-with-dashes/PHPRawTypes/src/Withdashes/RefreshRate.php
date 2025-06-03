<?php

namespace Grafana\Foundation\Withdashes;

/**
 * Refresh rate or disabled.
 */
class RefreshRate implements \JsonSerializable
{
    public ?string $string;

    public ?bool $bool;

    /**
     * @param string|null $string
     * @param bool|null $bool
     */
    public function __construct(?string $string = null, ?bool $bool = null)
    {
        $this->string = $string;
        $this->bool = $bool;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{String?: string, Bool?: bool} $inputData */
        $data = $inputData;
        return new self(
            string: $data["String"] ?? null,
            bool: $data["Bool"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->string)) {
            $data->String = $this->string;
        }
        if (isset($this->bool)) {
            $data->Bool = $this->bool;
        }
        return $data;
    }
}
