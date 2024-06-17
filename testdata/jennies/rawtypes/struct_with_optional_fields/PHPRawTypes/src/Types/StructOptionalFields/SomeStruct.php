<?php

namespace Grafana\Foundation\Types\StructOptionalFields;

class SomeStruct implements \JsonSerializable {
    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeOtherStruct $fieldRef;

    public ?string $fieldString;

    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeStructOperator $operator;

    /**
     * @var array<string>|null
     */
    public ?array $fieldArrayOfStrings;

    public ?\Grafana\Foundation\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

    /**
     * @param \Grafana\Foundation\Types\StructOptionalFields\SomeOtherStruct|null $fieldRef
     * @param string|null $fieldString
     * @param \Grafana\Foundation\Types\StructOptionalFields\SomeStructOperator|null $operator
     * @param array<string>|null $fieldArrayOfStrings
     * @param \Grafana\Foundation\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct|null $fieldAnonymousStruct
     */
    public function __construct(?\Grafana\Foundation\Types\StructOptionalFields\SomeOtherStruct $fieldRef = null, ?string $fieldString = null, ?\Grafana\Foundation\Types\StructOptionalFields\SomeStructOperator $operator = null, ?array $fieldArrayOfStrings = null, ?\Grafana\Foundation\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct = null)
    {
        $this->fieldRef = $fieldRef;
        $this->fieldString = $fieldString;
        $this->operator = $operator;
        $this->fieldArrayOfStrings = $fieldArrayOfStrings;
        $this->fieldAnonymousStruct = $fieldAnonymousStruct;
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        $data = [
        ];
        if (isset($this->fieldRef)) {
            $data["FieldRef"] = $this->fieldRef;
        }
        if (isset($this->fieldString)) {
            $data["FieldString"] = $this->fieldString;
        }
        if (isset($this->operator)) {
            $data["Operator"] = $this->operator;
        }
        if (isset($this->fieldArrayOfStrings)) {
            $data["FieldArrayOfStrings"] = $this->fieldArrayOfStrings;
        }
        if (isset($this->fieldAnonymousStruct)) {
            $data["FieldAnonymousStruct"] = $this->fieldAnonymousStruct;
        }
        return $data;
    }
}
