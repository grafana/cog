<?php

namespace Grafana\Foundation\Types\StructOptionalFields;

class SomeStruct implements \JsonSerializable {
    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeOtherStruct $fieldRef;

    public ?string $fieldString;

    public ?\Grafana\Foundation\Types\StructOptionalFields\SomeStructOperator $operator;

    /**
     * @var array<string>
     */
    public ?array $fieldArrayOfStrings;

    public ?\Grafana\Foundation\Types\StructOptionalFields\StructOptionalFieldsSomeStructFieldAnonymousStruct $fieldAnonymousStruct;

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
