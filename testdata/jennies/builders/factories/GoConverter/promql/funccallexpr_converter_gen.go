package promql



import (
	"strings"
	"fmt"
	cog "github.com/grafana/cog/generated/cog"
)

// FuncCallExprConverter accepts a `FuncCallExpr` object and generates the Go code to build this object using builders.
func FuncCallExprConverter(input FuncCallExpr) string {
    calls := []string{
    `promql.NewFuncCallExprBuilder()`,
    }
    var buffer strings.Builder
        if input.Function != "" {
        
    buffer.WriteString(`Function(`)
        arg0 :=fmt.Sprintf("%#v", input.Function)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.Args != nil && len(input.Args) >= 1 {
        
    buffer.WriteString(`Args(`)
        tmparg0 := []string{}
        for _, arg1 := range input.Args {
        tmpargsarg1 :=cog.Dump(arg1)
        tmparg0 = append(tmparg0, tmpargsarg1)
        }
        arg0 := "[]cog.Builder[promql.Expr]{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
