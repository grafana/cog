package sandbox



import (
	cog "github.com/grafana/cog/generated/cog"
)

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
