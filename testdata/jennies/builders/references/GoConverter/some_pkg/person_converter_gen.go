package some_pkg



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
)

// PersonConverter accepts a `Person` object and generates the Go code to build this object using builders.
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
