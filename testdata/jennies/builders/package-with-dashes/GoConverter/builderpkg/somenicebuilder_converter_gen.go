package builderpkg



import (
	"strings"
	withdashes "github.com/grafana/cog/generated/with-dashes"
	"fmt"
)

func SomeNiceBuilderConverter(input withdashes.SomeStruct) string {
    calls := []string{
    `builderpkg.NewSomeNiceBuilderBuilder()`,
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

    return strings.Join(calls, ".\t\n")
    }
