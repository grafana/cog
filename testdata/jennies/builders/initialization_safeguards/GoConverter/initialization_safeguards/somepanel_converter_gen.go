package initialization_safeguards



import (
	cog "github.com/grafana/cog/generated/cog"
)

func SomePanelConverter(input SomePanel) string {
    calls := []string{
    `initialization_safeguards.NewSomePanelBuilder()`,
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
        if input.Options != nil {
        
    buffer.WriteString(`ShowLegend(`)
        arg0 :=fmt.Sprintf("%#v", input.Options.Legend.Show)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
