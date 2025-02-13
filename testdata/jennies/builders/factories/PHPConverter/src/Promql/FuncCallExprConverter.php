<?php

namespace Grafana\Foundation\Promql;

final class FuncCallExprConverter
{
    public static function convert(\Grafana\Foundation\Promql\FuncCallExpr $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Promql\FuncCallExprBuilder())',
        ];
            if ($input->function !== "") {
    
        
    $buffer = 'function(';
        $arg0 =\var_export($input->function, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if (count($input->args) >= 1) {
    
        
    $buffer = 'args(';
        $tmparg0 = [];
        foreach ($input->args as $arg1) {
        $tmpargsarg1 =\var_export($arg1, true);
        $tmparg0[] = $tmpargsarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

