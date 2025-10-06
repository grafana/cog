<?php

namespace Grafana\Foundation\MapOfDisjunctions;

final class LibraryPanelConverter
{
    public static function convert(\Grafana\Foundation\MapOfDisjunctions\LibraryPanel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfDisjunctions\LibraryPanelBuilder())',
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
