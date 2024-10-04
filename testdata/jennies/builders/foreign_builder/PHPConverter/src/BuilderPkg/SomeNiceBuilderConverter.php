<?php

namespace Grafana\Foundation\BuilderPkg;

final class SomeNiceBuilderConverter
{
    public static function convert(\Grafana\Foundation\SomePkg\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\BuilderPkg\SomeNiceBuilderBuilder())',
        ];
            if ($input->title !== "") {
    
        
    $buffer = 'title(';
        $arg0 =\var_export($input->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

