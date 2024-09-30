package nullable_map_assignment



import (
	cog "github.com/grafana/cog/generated/cog"
)

func SomeStructConverter(input *SomeStruct) string {
    calls := []string{
    `nullable_map_assignment.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Config != nil {
        
    buffer.WriteString(`Config(`)
        arg0 :=cog.Dump(input.Config)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
