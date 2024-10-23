package sandbox



import (
	"strings"
	"fmt"
)

// DashboardConverter accepts a `Dashboard` object and generates the Go code to build this object using builders.
func DashboardConverter(input Dashboard) string {
    calls := []string{
    `sandbox.NewDashboardBuilder()`,
    }
    var buffer strings.Builder
        if input.Variables != nil && len(input.Variables) >= 1 {for _, item := range input.Variables {
        
    buffer.WriteString(`WithVariable(`)
        arg0 :=fmt.Sprintf("%#v", item.Name)
        buffer.WriteString(arg0)
        buffer.WriteString(", ")
        arg1 :=fmt.Sprintf("%#v", item.Value)
        buffer.WriteString(arg1)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }}

    return strings.Join(calls, ".\t\n")
}
