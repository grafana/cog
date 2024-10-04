<?php

namespace Grafana\Foundation\InitializationSafeguards;

final class SomePanelConverter
{
    public static function convert(\Grafana\Foundation\InitializationSafeguards\SomePanel $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\InitializationSafeguards\SomePanelBuilder())',
        ];
            if ($input->title !== "") {
    
        
    $buffer = 'title(';
        $arg0 =\var_export($input->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->options !== null) {
    
        
    $buffer = 'showLegend(';
        $arg0 =\var_export($input->options->legend->show, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

