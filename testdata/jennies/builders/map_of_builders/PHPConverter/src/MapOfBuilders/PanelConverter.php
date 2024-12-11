<?php

namespace Grafana\Foundation\MapOfBuilders;

final class PanelConverter
{
    public static function convert(\Grafana\Foundation\MapOfBuilders\Panel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfBuilders\PanelBuilder())',
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

