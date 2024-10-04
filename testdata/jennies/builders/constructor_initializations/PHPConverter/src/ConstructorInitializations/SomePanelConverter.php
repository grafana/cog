<?php

namespace Grafana\Foundation\ConstructorInitializations;

final class SomePanelConverter
{
    public static function convert(\Grafana\Foundation\ConstructorInitializations\SomePanel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\ConstructorInitializations\SomePanelBuilder())',
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

