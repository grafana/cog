package some_pkg



import (
	cog "github.com/grafana/cog/generated/cog"
)

func PersonConverter(input Person) string {
    calls := []string{
    `some_pkg.NewPersonBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Name(`)
        arg0 :=cog.Dump(input.Name)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
    }
