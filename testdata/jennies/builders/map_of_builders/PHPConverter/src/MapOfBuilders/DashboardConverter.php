<?php

namespace Grafana\Foundation\MapOfBuilders;

final class DashboardConverter
{
    public static function convert(\Grafana\Foundation\MapOfBuilders\Dashboard $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\MapOfBuilders\DashboardBuilder())',
        ];
            
    
        {
    $buffer = 'panels(';
        $arg0 = "[";
        foreach ($input->panels as $key => $arg1) {
            $tmppanelsarg1 = \Grafana\Foundation\MapOfBuilders\PanelConverter::convert($arg1);
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

