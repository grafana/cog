package map_of_disjunctions



import (
	"strings"
)

// ElementConverter accepts a `Element` object and generates the Go code to build this object using builders.
func ElementConverter(input Element) string {
    calls := []string{
    `map_of_disjunctions.NewElementBuilder()`,
    }
    var buffer strings.Builder
        if input.Panel != nil {
        
    buffer.WriteString(`Panel(`)
        arg0 := PanelConverter(*input.Panel)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.LibraryPanel != nil {
        
    buffer.WriteString(`LibraryPanel(`)
        arg0 := LibraryPanelConverter(*input.LibraryPanel)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
