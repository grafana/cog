package sandbox



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
)

// DashboardConverter accepts a `Dashboard` object and generates the Go code to build this object using builders.
func DashboardConverter(input Dashboard) string {
    calls := []string{
    `sandbox.NewDashboardBuilder()`,
    }
    var buffer strings.Builder
        if input.Variables != nil && len(input.Variables) >= 1 {
        
    buffer.WriteString(`Variables(`)
        tmparg0 := []string{}
        for _, arg1 := range input.Variables {
        tmpvariablesarg1 :=cog.Dump(arg1)
        tmparg0 = append(tmparg0, tmpvariablesarg1)
        }
        arg0 := "[]sandbox.Variable{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
