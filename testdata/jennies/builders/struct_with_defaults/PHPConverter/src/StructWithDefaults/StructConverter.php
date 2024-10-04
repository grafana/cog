<?php

namespace Grafana\Foundation\StructWithDefaults;

final class StructConverter
{
    public static function convert(\Grafana\Foundation\StructWithDefaults\Struct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\StructWithDefaults\StructBuilder())',
        ];
            
    
        {
    $buffer = 'allFields(';
        $arg0 = \Grafana\Foundation\StructWithDefaults\NestedStructConverter::convert($input->allFields);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            
    
        {
    $buffer = 'partialFields(';
        $arg0 = \Grafana\Foundation\StructWithDefaults\NestedStructConverter::convert($input->partialFields);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            
    
        {
    $buffer = 'emptyFields(';
        $arg0 = \Grafana\Foundation\StructWithDefaults\NestedStructConverter::convert($input->emptyFields);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            
    
        {
    $buffer = 'complexField(';
        $arg0 =\var_export($input->complexField, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            
    
        {
    $buffer = 'partialComplexField(';
        $arg0 =\var_export($input->partialComplexField, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

