<?php

namespace Grafana\Foundation\Constraints;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Constraints\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Constraints\SomeStructBuilder())',
        ];
        $buffer = '';
            
    
        {
    $buffer .= 'id(';
        $arg0 =\var_export($input->id, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    $buffer = '';
    }
    
    
            if ($input->title !== "") {
    
        
    $buffer .= 'title(';
        $arg0 =\var_export($input->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    $buffer = '';
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

