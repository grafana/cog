package builder_pkg



import (
	some_pkg "github.com/grafana/cog/generated/some_pkg"
	"strings"
	"fmt"
)

// SomeNiceBuilderConverter accepts a `SomeNiceBuilder` object and generates the Go code to build this object using builders.
func SomeNiceBuilderConverter(input some_pkg.SomeStruct) string {
    calls := []string{
    `builder_pkg.NewSomeNiceBuilderBuilder()`,
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
