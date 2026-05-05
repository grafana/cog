<?php

namespace Grafana\Foundation\Dashboard;

class DataSourceRef implements \JsonSerializable
{
    public ?string $type;

    public ?string $uid;

    /**
     * @param string|null $type
     * @param string|null $uid
     */
    public function __construct(?string $type = null, ?string $uid = null)
    {
        $this->type = $type;
        $this->uid = $uid;
    }

    public function equals(mixed $other): bool
    {
        if (!($other instanceof self)) {
            return false;
        }
    
        if (($this->type === null) !== ($other->type === null)) {
            return false;
        }
        if ($this->type !== null) {
            if ($this->type !== $other->type) {
                return false;
            }
        }
    
        if (($this->uid === null) !== ($other->uid === null)) {
            return false;
        }
        if ($this->uid !== null) {
            if ($this->uid !== $other->uid) {
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
        /** @var array{type?: string, uid?: string} $inputData */
        $data = $inputData;
        return new self(
            type: $data["type"] ?? null,
            uid: $data["uid"] ?? null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        if (isset($this->type)) {
            $data->type = $this->type;
        }
        if (isset($this->uid)) {
            $data->uid = $this->uid;
        }
        return $data;
    }
}
