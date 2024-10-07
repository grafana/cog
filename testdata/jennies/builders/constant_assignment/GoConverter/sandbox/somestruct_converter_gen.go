package sandbox



import (
	"strings"
)

func SomeStructConverter(input *SomeStruct) string {
    calls := []string{
    `sandbox.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Editable == true {
        
    buffer.WriteString(`Editable(`)
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.Editable == false {
        
    buffer.WriteString(`Readonly(`)
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.AutoRefresh != nil && *input.AutoRefresh == true {
        
    buffer.WriteString(`AutoRefresh(`)
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.AutoRefresh != nil && *input.AutoRefresh == false {
        
    buffer.WriteString(`NoAutoRefresh(`)
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
