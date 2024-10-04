<?php

namespace Grafana\Foundation\Properties;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Properties\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Properties\SomeStructBuilder())',
        ];
            
    
        {
    $buffer = 'id(';
        $arg0 =\var_export($input->id, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

