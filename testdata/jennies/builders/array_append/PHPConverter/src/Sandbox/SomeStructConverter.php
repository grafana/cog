<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructBuilder())',
        ];
            if (count($input->tags) >= 1) {
    foreach ($input->tags as $item) {
        
    $buffer = 'tags(';
        $arg0 =\var_export($item, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    }
    }

        return \implode("\n\t->", $calls);
    }
}

