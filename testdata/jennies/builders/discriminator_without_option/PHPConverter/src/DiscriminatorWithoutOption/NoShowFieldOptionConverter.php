<?php

namespace Grafana\Foundation\DiscriminatorWithoutOption;

final class NoShowFieldOptionConverter
{
    public static function convert(\Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOption $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\DiscriminatorWithoutOption\NoShowFieldOptionBuilder())',
        ];
            if ($input->text !== "") {
    
        
    $buffer = 'text(';
        $arg0 =\var_export($input->text, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}
