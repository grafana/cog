package sandbox



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
)

// SomeStructConverter accepts a `SomeStruct` object and generates the Go code to build this object using builders.
func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `sandbox.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Data != nil {for key, value := range input.Data {
        
    buffer.WriteString(`Data(`)
        arg0 :=cog.Dump(key)
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
