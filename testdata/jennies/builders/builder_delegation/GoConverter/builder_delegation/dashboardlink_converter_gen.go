package builder_delegation



import (
	"strings"
	"fmt"
)

// DashboardLinkConverter accepts a `DashboardLink` object and generates the Go code to build this object using builders.
func DashboardLinkConverter(input DashboardLink) string {
    calls := []string{
    `builder_delegation.NewDashboardLinkBuilder()`,
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
        if input.Url != "" {
        
    buffer.WriteString(`Url(`)
        arg0 :=fmt.Sprintf("%#v", input.Url)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
