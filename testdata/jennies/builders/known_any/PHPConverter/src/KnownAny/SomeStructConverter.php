<?php

namespace Grafana\Foundation\KnownAny;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\KnownAny\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\KnownAny\SomeStructBuilder())',
        ];
            if ($input->config !== null && $input->config instanceof \Grafana\Foundation\KnownAny\Config && $input->config->title !== "") {
    
        
    $buffer = 'title(';
        $arg0 =\var_export($input->config->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

