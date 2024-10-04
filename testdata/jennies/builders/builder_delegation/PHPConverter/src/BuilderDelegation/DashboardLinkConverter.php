<?php

namespace Grafana\Foundation\BuilderDelegation;

final class DashboardLinkConverter
{
    public static function convert(\Grafana\Foundation\BuilderDelegation\DashboardLink $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\BuilderDelegation\DashboardLinkBuilder())',
        ];
            if ($input->title !== "") {
    
        
    $buffer = 'title(';
        $arg0 =\var_export($input->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->url !== "") {
    
        
    $buffer = 'url(';
        $arg0 =\var_export($input->url, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

