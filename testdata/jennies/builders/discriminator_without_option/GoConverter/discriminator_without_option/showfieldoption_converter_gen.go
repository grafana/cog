package discriminator_without_option



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
)

// ShowFieldOptionConverter accepts a `ShowFieldOption` object and generates the Go code to build this object using builders.
func ShowFieldOptionConverter(input ShowFieldOption) string {
    calls := []string{
    `discriminator_without_option.NewShowFieldOptionBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Field(`)
        arg0 :=cog.Dump(input.Field)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        if input.Text != "" {
        
    buffer.WriteString(`Text(`)
        arg0 :=fmt.Sprintf("%#v", input.Text)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
