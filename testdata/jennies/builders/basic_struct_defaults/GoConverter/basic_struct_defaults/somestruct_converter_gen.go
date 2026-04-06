package basic_struct_defaults



import (
	"strings"
	"fmt"
)

// SomeStructConverter accepts a `SomeStruct` object and generates the Go code to build this object using builders.
func SomeStructConverter(input SomeStruct) string {
    calls := []string{
    `basic_struct_defaults.NewSomeStructBuilder()`,
    }
    var buffer strings.Builder
        if input.Id != 42 {
        
    buffer.WriteString(`Id(`)
        arg0 :=fmt.Sprintf("%#v", input.Id)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.Uid != "" && input.Uid != "default-uid" {
        
    buffer.WriteString(`Uid(`)
        arg0 :=fmt.Sprintf("%#v", input.Uid)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.Tags != nil && len(input.Tags) >= 1 {
        
    buffer.WriteString(`Tags(`)
        tmparg0 := []string{}
        for _, arg1 := range input.Tags {
        tmptagsarg1 :=fmt.Sprintf("%#v", arg1)
        tmparg0 = append(tmparg0, tmptagsarg1)
        }
        arg0 := "[]string{" + strings.Join(tmparg0, ",\n") + "}"
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }
        if input.LiveNow != true {
        
    buffer.WriteString(`LiveNow(`)
        arg0 :=fmt.Sprintf("%#v", input.LiveNow)
        buffer.WriteString(arg0)
        
    buffer.WriteString(")")

    calls = append(calls, buffer.String())
    buffer.Reset()
    
    }

    return strings.Join(calls, ".\t\n")
}
