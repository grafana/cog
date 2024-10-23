package struct_with_defaults



import (
	"strings"
	cog "github.com/grafana/cog/generated/cog"
)

// StructConverter accepts a `Struct` object and generates the Go code to build this object using builders.
func StructConverter(input Struct) string {
    calls := []string{
    `struct_with_defaults.NewStructBuilder()`,
    }
    var buffer strings.Builder
        
        {
    buffer.WriteString(`AllFields(`)
        arg0 := NestedStructConverter(input.AllFields)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        
        {
    buffer.WriteString(`PartialFields(`)
        arg0 := NestedStructConverter(input.PartialFields)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        
        {
    buffer.WriteString(`EmptyFields(`)
        arg0 := NestedStructConverter(input.EmptyFields)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        
        {
    buffer.WriteString(`ComplexField(`)
        arg0 :=cog.Dump(input.ComplexField)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    
        
        {
    buffer.WriteString(`PartialComplexField(`)
        arg0 :=cog.Dump(input.PartialComplexField)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    }
    

    return strings.Join(calls, ".\t\n")
}
