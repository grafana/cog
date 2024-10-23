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
        if input.Time != nil && input.Time.From != "" && input.Time.From != "now-6h" && input.Time.To != "" && input.Time.To != "now" {
        
    buffer.WriteString(`Time(`)
        arg0 :=fmt.Sprintf("%#v", input.Time.From)
        buffer.WriteString(arg0)
        buffer.WriteString(", ")
        arg1 :=fmt.Sprintf("%#v", input.Time.To)
        buffer.WriteString(arg1)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
