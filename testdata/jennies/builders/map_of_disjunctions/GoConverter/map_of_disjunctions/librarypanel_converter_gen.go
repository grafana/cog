package map_of_disjunctions



import (
	"strings"
	"fmt"
)

// LibraryPanelConverter accepts a `LibraryPanel` object and generates the Go code to build this object using builders.
func LibraryPanelConverter(input LibraryPanel) string {
    calls := []string{
    `map_of_disjunctions.NewLibraryPanelBuilder()`,
    }
    var buffer strings.Builder
        if input.Text != "" {
        
    buffer.WriteString(`Text(`)
        arg0 :=fmt.Sprintf("%#v", input.Text)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
