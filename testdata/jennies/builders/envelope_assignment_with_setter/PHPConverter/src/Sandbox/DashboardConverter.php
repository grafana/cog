<?php

namespace Grafana\Foundation\Sandbox;

final class DashboardConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\Dashboard $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\DashboardBuilder())',
        ];
            if (count($input->variables) >= 1) {
    
        
    $buffer = 'variables(';
        $tmparg0 = [];
        foreach ($input->variables as $arg1) {
        $tmpvariablesarg1 ='(new \Grafana\Foundation\Sandbox\Variable(name: '.\var_export($arg1->name, true).',value: '.\var_export($arg1->value, true).',))';
        $tmparg0[] = $tmpvariablesarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

