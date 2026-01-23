<?php

namespace Grafana\Foundation\MapOfDisjunctions;

final class PanelOrLibraryPanelConverter
{
    public static function convert(\Grafana\Foundation\MapOfDisjunctions\PanelOrLibraryPanel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfDisjunctions\PanelOrLibraryPanelBuilder())',
        ];
            if ($input->panel !== null) {
    
        
    $buffer = 'panel(';
        $arg0 = \Grafana\Foundation\MapOfDisjunctions\PanelConverter::convert($input->panel);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->libraryPanel !== null) {
    
        
    $buffer = 'libraryPanel(';
        $arg0 = \Grafana\Foundation\MapOfDisjunctions\LibraryPanelConverter::convert($input->libraryPanel);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}
