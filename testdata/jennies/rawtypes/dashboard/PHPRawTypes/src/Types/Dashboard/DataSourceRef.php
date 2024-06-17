<?php

namespace Grafana\Foundation\Types\Dashboard;

class DataSourceRef implements \JsonSerializable {
    public ?string $type;

    public ?string $uid;

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
