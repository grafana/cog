package struct_with_defaults



import (
	"strings"
	"fmt"
)

// NestedStructConverter accepts a `NestedStruct` object and generates the Go code to build this object using builders.
func NestedStructConverter(input NestedStruct) string {
    calls := []string{
    `struct_with_defaults.NewNestedStructBuilder()`,
    }
    var buffer strings.Builder
        if input.StringVal != "" {
        
    buffer.WriteString(`StringVal(`)
        arg0 :=fmt.Sprintf("%#v", input.StringVal)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        
        {
    buffer.WriteString(`IntVal(`)
        arg0 :=fmt.Sprintf("%#v", input.IntVal)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
}
