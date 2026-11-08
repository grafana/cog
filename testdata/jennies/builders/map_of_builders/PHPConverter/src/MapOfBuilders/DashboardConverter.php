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
    
    
            
    
        {
    $buffer = 'facets(';
        $arg0 = "[";
        foreach ($input->facets as $key => $arg1) {
            $tmptmpfacetsarg1 = [];
        foreach ($arg1 as $arg1Value) {
        $tmparg1arg1Value = \Grafana\Foundation\MapOfBuilders\FacetValueConverter::convert($arg1Value);
        $tmptmpfacetsarg1[] = $tmparg1arg1Value;
        }
        $tmpfacetsarg1 = "[" . implode(", \n", $tmptmpfacetsarg1) . "]";
            $arg0 .= "\t".var_export($key, true)." => $tmpfacetsarg1,";
        }
        $arg0 .= "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

