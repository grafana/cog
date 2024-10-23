package known_any



import (
	"strings"
	"fmt"
)

// SomeStructConverter accepts a `SomeStruct` object and generates the Go code to build this object using builders.
func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `known_any.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Config != nil && input.Config.(*Config).Title != "" {
        
    buffer.WriteString(`Title(`)
        arg0 :=fmt.Sprintf("%#v", input.Config.(*Config).Title)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
