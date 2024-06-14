<?php

namespace Types\Defaults;

class Struct
{
    public \Types\Defaults\NestedStruct $allFields;

    public \Types\Defaults\NestedStruct $partialFields;

    public \Types\Defaults\NestedStruct $emptyFields;

    public \Types\Defaults\DefaultsStructComplexField $complexField;

    public \Types\Defaults\DefaultsStructPartialComplexField $partialComplexField;
}
