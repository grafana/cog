<?php

namespace Grafana\Foundation\StructWithDefaults;

final class NestedStructConverter
{
    public static function convert(\Grafana\Foundation\StructWithDefaults\NestedStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\StructWithDefaults\NestedStructBuilder())',
        ];
        $buffer = '';
            if ($input->stringVal !== "") {
    
        
    $buffer .= 'stringVal(';
        $arg0 =\var_export($input->stringVal, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    $buffer = '';
    
    
    }
            
    
        {
    $buffer .= 'intVal(';
        $arg0 =\var_export($input->intVal, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    $buffer = '';
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

