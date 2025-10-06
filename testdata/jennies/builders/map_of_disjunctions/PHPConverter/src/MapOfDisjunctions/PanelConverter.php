<?php

namespace Grafana\Foundation\MapOfDisjunctions;

final class PanelConverter
{
    public static function convert(\Grafana\Foundation\MapOfDisjunctions\Panel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfDisjunctions\PanelBuilder())',
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
