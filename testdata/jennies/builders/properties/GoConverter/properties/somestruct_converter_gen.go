package properties



import (
	cog "github.com/grafana/cog/generated/cog"
)

func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `properties.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Id(`)
        arg0 :=fmt.Sprintf("%#v", input.Id)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
    }
