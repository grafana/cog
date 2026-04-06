package composable_slot



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
)

// LokiBuilderConverter accepts a `LokiBuilder` object and generates the Go code to build this object using builders.
func LokiBuilderConverter(input Dashboard) string {
    calls := []string{
    `composable_slot.NewLokiBuilderBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Target(`)
        arg0 := cog.ConvertDataqueryToCode(input.Target, )
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        if input.Targets != nil && len(input.Targets) >= 1 {
        
    buffer.WriteString(`Targets(`)
        tmparg0 := []string{}
        for _, arg1 := range input.Targets {
        tmptargetsarg1 := cog.ConvertDataqueryToCode(arg1, )
        tmparg0 = append(tmparg0, tmptargetsarg1)
        }
        arg0 := "[]cog.Builder[variants.Dataquery]{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
