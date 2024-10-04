package known_any



import (
	cog "github.com/grafana/cog/generated/cog"
)

func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `known_any.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Config != nil && input.Config.(*Config).Title != "" {
        
    buffer.WriteString(`Title(`)
        arg0 :=fmt.Sprintf("%#v", input.Config.(*Config).Title)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
