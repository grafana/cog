package dataquery_variant_builder



import (
	strings "strings"
	fmt "fmt"
)

func LokiBuilderConverter(input Loki) string {
    calls := []string{
    `dataquery_variant_builder.NewLokiBuilderBuilder()`,
    }
    var buffer strings.Builder
        if input.Expr != "" {
        
    buffer.WriteString(`Expr(`)
        arg0 :=fmt.Sprintf("%#v", input.Expr)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
    }
