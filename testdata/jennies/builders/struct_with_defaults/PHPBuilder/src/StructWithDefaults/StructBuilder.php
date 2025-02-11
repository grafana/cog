<?php

namespace Grafana\Foundation\StructWithDefaults;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\StructWithDefaults\Struct>
 */
class StructBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\StructWithDefaults\Struct $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\StructWithDefaults\Struct();
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\StructWithDefaults\Struct
     */
    public function build()
    {
        return $this->internal;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\StructWithDefaults\NestedStruct> $allFields
     */
    public function allFields(\Grafana\Foundation\Cog\Builder $allFields): static
    {
        $allFieldsResource = $allFields->build();
        $this->internal->allFields = $allFieldsResource;
    
        return $this;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\StructWithDefaults\NestedStruct> $partialFields
     */
    public function partialFields(\Grafana\Foundation\Cog\Builder $partialFields): static
    {
        $partialFieldsResource = $partialFields->build();
        $this->internal->partialFields = $partialFieldsResource;
    
        return $this;
    }

    /**
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\StructWithDefaults\NestedStruct> $emptyFields
     */
    public function emptyFields(\Grafana\Foundation\Cog\Builder $emptyFields): static
    {
        $emptyFieldsResource = $emptyFields->build();
        $this->internal->emptyFields = $emptyFieldsResource;
    
        return $this;
    }

    public function complexField(unknown $complexField): static
    {
        $this->internal->complexField = $complexField;
    
        return $this;
    }

    public function partialComplexField(unknown $partialComplexField): static
    {
        $this->internal->partialComplexField = $partialComplexField;
    
        return $this;
    }

}
