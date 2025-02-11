package struct_optional_fields;

import java.util.List;

public class SomeStruct {
    public SomeOtherStruct fieldRef;
    public String fieldString;
    public SomeStructOperator operator;
    public List<String> fieldArrayOfStrings;
    public StructOptionalFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct;
    public SomeStruct() {
    }
    public SomeStruct(SomeOtherStruct fieldRef,String fieldString,SomeStructOperator operator,List<String> fieldArrayOfStrings,StructOptionalFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct) {
        this.fieldRef = fieldRef;
        this.fieldString = fieldString;
        this.operator = operator;
        this.fieldArrayOfStrings = fieldArrayOfStrings;
        this.fieldAnonymousStruct = fieldAnonymousStruct;
    }
}
