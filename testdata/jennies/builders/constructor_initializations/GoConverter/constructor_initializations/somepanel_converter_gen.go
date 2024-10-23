package constructor_initializations



import (
	"strings"
	"fmt"
)

// SomePanelConverter accepts a `SomePanel` object and generates the Go code to build this object using builders.
func SomePanelConverter(input SomePanel) string {
    calls := []string{
    `constructor_initializations.NewSomePanelBuilder()`,
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
