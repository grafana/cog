<?php

namespace Grafana\Foundation\SelfReferentialStruct;

/**
 * Node represents a node in a singly-linked list.
 * The next field points to the following node, or is absent if this is the last node.
 */
class Node implements \JsonSerializable
{
    public string $value;

    public ?\Grafana\Foundation\SelfReferentialStruct\Node $next;

    /**
     * @param string|null $value
     * @param \Grafana\Foundation\SelfReferentialStruct\Node|null $next
     */
    public function __construct(?string $value = null, ?\Grafana\Foundation\SelfReferentialStruct\Node $next = null)
    {
        $this->value = $value ?: "";
        $this->next = $next;
    }

    /**
     * @param array<string, mixed> $inputData
     */
    public static function fromArray(array $inputData): self
    {
        /** @var array{value?: string, next?: mixed} $inputData */
        $data = $inputData;
        return new self(
            value: $data["value"] ?? null,
            next: isset($data["next"]) ? (function($input) {
    	/** @var array{value?: string, next?: mixed} */
    $val = $input;
    	return \Grafana\Foundation\SelfReferentialStruct\Node::fromArray($val);
    })($data["next"]) : null,
        );
    }

    /**
     * @return mixed
     */
    public function jsonSerialize(): mixed
    {
        $data = new \stdClass;
        $data->value = $this->value;
        if (isset($this->next)) {
            $data->next = $this->next;
        }
        return $data;
    }
}
