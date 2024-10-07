package sandbox



import (
	strings "strings"
	fmt "fmt"
)

func SomeStructConverter(input *SomeStruct) string {
    calls := []string{
    `sandbox.NewSomeStructBuilder(`+fmt.Sprintf("%#v", input.Title)+`)`,
    }
    var buffer strings.Builder
        if input.Title != "" {
        
    buffer.WriteString(`Title(`)
        arg0 :=fmt.Sprintf("%#v", input.Title)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
