package constraints



import (
	"strings"
	"fmt"
)

// SomeStructConverter accepts a `SomeStruct` object and generates the Go code to build this object using builders.
func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `constraints.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Id(`)
        arg0 :=fmt.Sprintf("%#v", input.Id)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
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
