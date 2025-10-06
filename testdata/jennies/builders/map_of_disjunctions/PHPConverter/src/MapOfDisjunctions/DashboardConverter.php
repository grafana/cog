<?php

namespace Grafana\Foundation\MapOfDisjunctions;

final class DashboardConverter
{
    public static function convert(\Grafana\Foundation\MapOfDisjunctions\Dashboard $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfDisjunctions\DashboardBuilder())',
        ];
            
    
        {
    $buffer = 'panels(';
        $arg0 = "[";
        foreach ($input->panels as $key => $arg1) {
            $tmppanelsarg1 = \Grafana\Foundation\MapOfDisjunctions\ElementConverter::convert($arg1);
            $arg0 .= "\t".var_export($key, true)." => $tmppanelsarg1,";
        }
        $arg0 .= "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}
