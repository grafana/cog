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
        if input.Tags != nil && len(input.Tags) >= 1 {for _, item := range input.Tags {
        
    buffer.WriteString(`Tags(`)
        arg0 :=fmt.Sprintf("%#v", item)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }}

    return strings.Join(calls, ".\t\n")
}
