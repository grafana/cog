package map_of_builders



import (
	"strings"
	"fmt"
)

// FacetValueConverter accepts a `FacetValue` object and generates the Go code to build this object using builders.
func FacetValueConverter(input FacetValue) string {
    calls := []string{
    `map_of_builders.NewFacetValueBuilder()`,
    }
    var buffer strings.Builder
        if input.Value != "" {
        
    buffer.WriteString(`Value(`)
        arg0 :=fmt.Sprintf("%#v", input.Value)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        
        {
    buffer.WriteString(`Count(`)
        arg0 :=fmt.Sprintf("%#v", input.Count)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
}
