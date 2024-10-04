<?php

namespace Grafana\Foundation\BasicStructDefaults;

final class SomeStructConverter
{
    public static function convert(\Grafana\Foundation\BasicStructDefaults\SomeStruct $input): string
    {
        
        $calls = [
            '(new \Grafana\Foundation\BasicStructDefaults\SomeStructBuilder())',
        ];
            if ($input->id !== 42) {
    
        
    $buffer = 'id(';
        $arg0 =\var_export($input->id, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->uid !== "" && $input->uid !== "default-uid") {
    
        
    $buffer = 'uid(';
        $arg0 =\var_export($input->uid, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if (count($input->tags) >= 1) {
    
        
    $buffer = 'tags(';
        $tmparg0 = [];
        foreach ($input->tags as $arg1) {
        $tmptagsarg1 =\var_export($arg1, true);
        $tmparg0[] = $tmptagsarg1;
        }
        $arg0 = "[" . implode(", \n", $tmparg0) . "]";
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }
            if ($input->liveNow !== true) {
    
        
    $buffer = 'liveNow(';
        $arg0 =\var_export($input->liveNow, true);
        $buffer .= $arg0;
        
    $buffer .= ')';

    $calls[] = $buffer;
    
    
    }

        return \implode("\n\t->", $calls);
    }
}

