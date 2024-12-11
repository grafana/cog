package map_of_builders



import (
	"strings"
	"fmt"
)

// DashboardConverter accepts a `Dashboard` object and generates the Go code to build this object using builders.
func DashboardConverter(input Dashboard) string {
    calls := []string{
    `map_of_builders.NewDashboardBuilder()`,
    }
    var buffer strings.Builder
        if input.Panels != nil {
        
    buffer.WriteString(`Panels(`)
        arg0 := "map[string]cog.Builder[map_of_builders.Panel]{"
        for key, arg1 := range input.Panels {
            tmppanelsarg1 := PanelConverter(arg1)
            arg0 += "\t" + fmt.Sprintf("%#v", key) + ": " + tmppanelsarg1 +","
        }
        arg0 += "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
