package nullable_map_assignment



import (
	cog "github.com/grafana/cog/generated/cog"
)

func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `nullable_map_assignment.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Config != nil {
        
    buffer.WriteString(`Config(`)
        arg0 := "map[string]string{"
        for key, arg1 := range input.Config {
            tmpconfigarg1 :=fmt.Sprintf("%#v", arg1)
            arg0 += "\t" + fmt.Sprintf("%#v", key) + ": " + tmpconfigarg1 +","
        }
        arg0 += "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
