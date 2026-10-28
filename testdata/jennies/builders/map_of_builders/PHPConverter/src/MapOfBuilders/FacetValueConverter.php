<?php

namespace Grafana\Foundation\MapOfBuilders;

final class FacetValueConverter
{
    public static function convert(\Grafana\Foundation\MapOfBuilders\FacetValue $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfBuilders\FacetValueBuilder())',
        ];
            if ($input->value !== "") {
    
        
    $buffer = 'value(';
        $arg0 =\var_export($input->value, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            
    
        {
    $buffer = 'count(';
        $arg0 =\var_export($input->count, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

