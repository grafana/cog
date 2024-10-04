<?php

namespace Grafana\Foundation\BuilderDelegation;

final class DashboardConverter
{
    public static function convert(\Grafana\Foundation\BuilderDelegation\Dashboard $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\BuilderDelegation\DashboardBuilder())',
        ];
            
    
        {
    $buffer = 'id(';
        $arg0 =\var_export($input->id, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    
            if ($input->title !== "") {
    
        
    $buffer = 'title(';
        $arg0 =\var_export($input->title, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if (count($input->links) >= 1) {
    
        
    $buffer = 'links(';
        $tmparg0 = [];
        foreach ($input->links as $arg1) {
        $tmplinksarg1 = \Grafana\Foundation\BuilderDelegation\DashboardLinkConverter::convert($arg1);
        $tmparg0[] = $tmplinksarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if (count($input->linksOfLinks) >= 1) {
    
        
    $buffer = 'linksOfLinks(';
        $tmparg0 = [];
        foreach ($input->linksOfLinks as $arg1) {
        $tmptmplinksOfLinksarg1 = [];
        foreach ($arg1 as $arg1Value) {
        $tmparg1arg1Value = \Grafana\Foundation\BuilderDelegation\DashboardLinkConverter::convert($arg1Value);
        $tmptmplinksOfLinksarg1[] = $tmparg1arg1Value;
        }
        $tmplinksOfLinksarg1 = "[" . implode(", \n", $tmptmplinksOfLinksarg1) . "]";
        $tmparg0[] = $tmplinksOfLinksarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            
    
        {
    $buffer = 'singleLink(';
        $arg0 = \Grafana\Foundation\BuilderDelegation\DashboardLinkConverter::convert($input->singleLink);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    }
    
    

        return \implode("\n\t->", $calls);
    }
}

