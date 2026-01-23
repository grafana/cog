package discriminator_without_option



import (
	"strings"
	"fmt"
)

// NoShowFieldOptionConverter accepts a `NoShowFieldOption` object and generates the Go code to build this object using builders.
func NoShowFieldOptionConverter(input NoShowFieldOption) string {
    calls := []string{
    `discriminator_without_option.NewNoShowFieldOptionBuilder()`,
    }
    var buffer strings.Builder
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
