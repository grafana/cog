package sandbox



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
)

// SomeStructWithDefaultEnumConverter accepts a `SomeStructWithDefaultEnum` object and generates the Go code to build this object using builders.
func SomeStructWithDefaultEnumConverter(input SomeStructWithDefaultEnum) string {
    calls := []string{
    `sandbox.NewSomeStructWithDefaultEnumBuilder()`,
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
