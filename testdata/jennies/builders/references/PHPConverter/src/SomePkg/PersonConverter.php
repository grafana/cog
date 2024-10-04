<?php

namespace Grafana\Foundation\SomePkg;

final class PersonConverter
{
    public static function convert(\Grafana\Foundation\SomePkg\Person $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\SomePkg\PersonBuilder())',
        ];
            
    
        {
    $buffer = 'name(';
        $arg0 =\var_export($input->name, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

