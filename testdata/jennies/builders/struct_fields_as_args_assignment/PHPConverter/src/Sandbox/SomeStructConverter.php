<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructBuilder())',
        ];
            if ($input->time !== null && $input->time->from !== "" && $input->time->from !== "now-6h" && $input->time->to !== "" && $input->time->to !== "now") {
    
        
    $buffer = 'time(';
        $arg0 =\var_export($input->time->from, true);
        $buffer .= $arg0;
        $buffer .= ', ';
        $arg1 =\var_export($input->time->to, true);
        $buffer .= $arg1;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

