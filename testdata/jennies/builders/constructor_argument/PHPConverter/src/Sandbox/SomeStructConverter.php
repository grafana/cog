<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStruct $input): string
    {
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructBuilder('.\var_export($input->title, true).'))',
        ];
        $buffer = '';
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

