<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructWithDefaultEnumConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStructWithDefaultEnum $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructWithDefaultEnumBuilder())',
        ];
            
    foreach ($input->data as $key => $value) {
        {
    $buffer = 'data(';
        $arg0 ='\Grafana\Foundation\Sandbox\StringEnumWithDefault::fromValue("'.$key.'")';
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

