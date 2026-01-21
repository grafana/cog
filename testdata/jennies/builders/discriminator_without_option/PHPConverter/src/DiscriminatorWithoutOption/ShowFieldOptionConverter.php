<?php

namespace Grafana\Foundation\DiscriminatorWithoutOption;

final class ShowFieldOptionConverter
{
    public static function convert(\Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOption $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\DiscriminatorWithoutOption\ShowFieldOptionBuilder())',
        ];
            
    
        {
    $buffer = 'field(';
        $arg0 ='\Grafana\Foundation\DiscriminatorWithoutOption\AnEnum::fromValue("'.$input->field.'")';
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
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
