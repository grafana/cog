package sandbox



import (
	"strings"
	"fmt"
)

// SomeStructConverter accepts a `SomeStruct` object and generates the Go code to build this object using builders.
func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `sandbox.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Annotations != nil {for key, value := range input.Annotations {
        
    buffer.WriteString(`Annotations(`)
        arg0 :=fmt.Sprintf("%#v", key)
        buffer.WriteString(arg0)
        buffer.WriteString(", ")
        arg1 :=fmt.Sprintf("%#v", value)
        buffer.WriteString(arg1)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }}

    return strings.Join(calls, ".\t\n")
}
