<?php

namespace Grafana\Foundation\Sandbox;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\Sandbox\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\Sandbox\SomeStructBuilder())',
        ];
            if ($input->editable === true) {
    
        
    $buffer = 'editable(';
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->editable === false) {
    
        
    $buffer = 'readonly(';
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->autoRefresh !== null && $input->autoRefresh === true) {
    
        
    $buffer = 'autoRefresh(';
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->autoRefresh !== null && $input->autoRefresh === false) {
    
        
    $buffer = 'noAutoRefresh(';
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

