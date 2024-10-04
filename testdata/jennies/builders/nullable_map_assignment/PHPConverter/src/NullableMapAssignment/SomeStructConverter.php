<?php

namespace Grafana\Foundation\NullableMapAssignment;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\NullableMapAssignment\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\NullableMapAssignment\SomeStructBuilder())',
        ];
            if ($input->config !== null) {
    
        
    $buffer = 'config(';
        $arg0 = "[";
        foreach ($input->config as $key => $arg1) {
            $tmpconfigarg1 =\var_export($arg1, true);
            $arg0 .= "\t".var_export($key, true)." => $tmpconfigarg1,";
        }
        $arg0 .= "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

