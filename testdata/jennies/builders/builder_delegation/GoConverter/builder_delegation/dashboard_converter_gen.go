package builder_delegation



import (
	"strings"
	"fmt"
)

// DashboardConverter accepts a `Dashboard` object and generates the Go code to build this object using builders.
func DashboardConverter(input Dashboard) string {
    calls := []string{
    `builder_delegation.NewDashboardBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`Id(`)
        arg0 :=fmt.Sprintf("%#v", input.Id)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        if input.Title != "" {
        
    buffer.WriteString(`Title(`)
        arg0 :=fmt.Sprintf("%#v", input.Title)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.Links != nil && len(input.Links) >= 1 {
        
    buffer.WriteString(`Links(`)
        tmparg0 := []string{}
        for _, arg1 := range input.Links {
        tmplinksarg1 := DashboardLinkConverter(arg1)
        tmparg0 = append(tmparg0, tmplinksarg1)
        }
        arg0 := "[]cog.Builder[builder_delegation.DashboardLink]{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.LinksOfLinks != nil && len(input.LinksOfLinks) >= 1 {
        
    buffer.WriteString(`LinksOfLinks(`)
        tmparg0 := []string{}
        for _, arg1 := range input.LinksOfLinks {
        tmptmplinksOfLinksarg1 := []string{}
        for _, arg1Value := range arg1 {
        tmparg1arg1Value := DashboardLinkConverter(arg1Value)
        tmptmplinksOfLinksarg1 = append(tmptmplinksOfLinksarg1, tmparg1arg1Value)
        }
        tmplinksOfLinksarg1 := "[]cog.Builder[builder_delegation.DashboardLink]{" + strings.Join(tmptmplinksOfLinksarg1, ",\n") + "}"
        tmparg0 = append(tmparg0, tmplinksOfLinksarg1)
        }
        arg0 := "[][]cog.Builder[builder_delegation.DashboardLink]{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        
        {
    buffer.WriteString(`SingleLink(`)
        arg0 := DashboardLinkConverter(input.SingleLink)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
}
