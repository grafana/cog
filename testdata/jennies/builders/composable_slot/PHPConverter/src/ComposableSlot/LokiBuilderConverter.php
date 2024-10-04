<?php

namespace Grafana\Foundation\ComposableSlot;

final class LokiBuilderConverter
{
    public static function convert(\Grafana\Foundation\ComposableSlot\Dashboard $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\ComposableSlot\LokiBuilderBuilder())',
        ];
            
    
        {
    $buffer = 'target(';
        $arg0 = \Grafana\Foundation\Cog\Runtime::get()->convertDataqueryToCode($input->target, );
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            if (count($input->targets) >= 1) {
    
        
    $buffer = 'targets(';
        $tmparg0 = [];
        foreach ($input->targets as $arg1) {
        $tmptargetsarg1 = \Grafana\Foundation\Cog\Runtime::get()->convertDataqueryToCode($arg1, );
        $tmparg0[] = $tmptargetsarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

