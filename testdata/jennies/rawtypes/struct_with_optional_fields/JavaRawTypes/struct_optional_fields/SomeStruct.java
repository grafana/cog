package struct_optional_fields;

import java.util.Objects;
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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldRef, o.fieldRef)) return false;
        if (!Objects.equals(this.fieldString, o.fieldString)) return false;
        if (!Objects.equals(this.operator, o.operator)) return false;
        if (!Objects.equals(this.fieldArrayOfStrings, o.fieldArrayOfStrings)) return false;
        if (!Objects.equals(this.fieldAnonymousStruct, o.fieldAnonymousStruct)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldRef, this.fieldString, this.operator, this.fieldArrayOfStrings, this.fieldAnonymousStruct);
    }
}
