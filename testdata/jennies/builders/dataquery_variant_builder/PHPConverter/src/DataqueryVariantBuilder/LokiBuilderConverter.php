<?php

namespace Grafana\Foundation\DataqueryVariantBuilder;

final class LokiBuilderConverter
{
    public static function convert(\Grafana\Foundation\DataqueryVariantBuilder\Loki $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\DataqueryVariantBuilder\LokiBuilderBuilder())',
        ];
            if ($input->expr !== "") {
    
        
    $buffer = 'expr(';
        $arg0 =\var_export($input->expr, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

