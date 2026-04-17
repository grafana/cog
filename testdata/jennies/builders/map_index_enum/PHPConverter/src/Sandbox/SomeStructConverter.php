<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructBuilder())',
        ];
            
    foreach ($input->data as $key => $value) {
        {
    $buffer = 'data(';
        $arg0 ='\Grafana\Foundation\Sandbox\StringEnum::fromValue("'.$key.'")';
        $buffer .= $arg0;
        $buffer .= ', ';
        $arg1 =\var_export($value, true);
        $buffer .= $arg1;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    }
    

        return \implode("\n\t->", $calls);
    }
}

