package map_of_builders



import (
	"strings"
	"fmt"
)

// PanelConverter accepts a `Panel` object and generates the Go code to build this object using builders.
func PanelConverter(input Panel) string {
    calls := []string{
    `map_of_builders.NewPanelBuilder()`,
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
