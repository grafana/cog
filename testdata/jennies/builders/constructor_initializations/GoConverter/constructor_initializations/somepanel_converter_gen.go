package constructor_initializations



import (
	strings "strings"
	fmt "fmt"
)

func SomePanelConverter(input *SomePanel) string {
    calls := []string{
    `constructor_initializations.NewSomePanelBuilder()`,
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
