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
    foreach ($input->variables as $item) {
        
    $buffer = 'withVariable(';
        $arg0 =\var_export($item->name, true);
        $buffer .= $arg0;
        $buffer .= ', ';
        $arg1 =\var_export($item->value, true);
        $buffer .= $arg1;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    }
    }

        return \implode("\n\t->", $calls);
    }
}

