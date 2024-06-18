<?php

namespace Grafana\Foundation\Dashboard;

class DataSourceRef implements \JsonSerializable {
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

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
        ];
        if (isset($this->type)) {
            $data["type"] = $this->type;
        }
        if (isset($this->uid)) {
            $data["uid"] = $this->uid;
        }
        return $data;
    }
}
