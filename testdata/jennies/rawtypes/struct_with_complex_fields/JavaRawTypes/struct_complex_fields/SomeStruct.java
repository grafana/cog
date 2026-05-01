package struct_complex_fields;

import java.util.Objects;
import java.util.LinkedList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

// This struct does things.
public class SomeStruct {
    public SomeOtherStruct fieldRef;
    public StringOrBool fieldDisjunctionOfScalars;
    public StringOrSomeOtherStruct fieldMixedDisjunction;
    public String fieldDisjunctionWithNull;
    public SomeStructOperator operator;
    public List<String> fieldArrayOfStrings;
    public Map<String, String> fieldMapOfStringToString;
    public StructComplexFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct;
    public String fieldRefToConstant;
    public SomeStruct() {
        this.fieldRef = new struct_complex_fields.SomeOtherStruct();
        this.fieldDisjunctionOfScalars = new struct_complex_fields.StringOrBool();
        this.fieldMixedDisjunction = new struct_complex_fields.StringOrSomeOtherStruct();
        this.operator = SomeStructOperator.GREATER_THAN;
        this.fieldArrayOfStrings = new LinkedList<>();
        this.fieldMapOfStringToString = new HashMap<>();
        this.fieldAnonymousStruct = new struct_complex_fields.StructComplexFieldsSomeStructFieldAnonymousStruct();
        this.fieldRefToConstant = new struct_complex_fields.ConnectionPath();
    }
    public SomeStruct(SomeOtherStruct fieldRef,StringOrBool fieldDisjunctionOfScalars,StringOrSomeOtherStruct fieldMixedDisjunction,String fieldDisjunctionWithNull,SomeStructOperator operator,List<String> fieldArrayOfStrings,Map<String, String> fieldMapOfStringToString,StructComplexFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct,String fieldRefToConstant) {
        this.fieldRef = fieldRef;
        this.fieldDisjunctionOfScalars = fieldDisjunctionOfScalars;
        this.fieldMixedDisjunction = fieldMixedDisjunction;
        this.fieldDisjunctionWithNull = fieldDisjunctionWithNull;
        this.operator = operator;
        this.fieldArrayOfStrings = fieldArrayOfStrings;
        this.fieldMapOfStringToString = fieldMapOfStringToString;
        this.fieldAnonymousStruct = fieldAnonymousStruct;
        this.fieldRefToConstant = fieldRefToConstant;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldRef, o.fieldRef)) return false;
        if (!Objects.equals(this.fieldDisjunctionOfScalars, o.fieldDisjunctionOfScalars)) return false;
        if (!Objects.equals(this.fieldMixedDisjunction, o.fieldMixedDisjunction)) return false;
        if (!Objects.equals(this.fieldDisjunctionWithNull, o.fieldDisjunctionWithNull)) return false;
        if (!Objects.equals(this.operator, o.operator)) return false;
        if (!Objects.equals(this.fieldArrayOfStrings, o.fieldArrayOfStrings)) return false;
        if (!Objects.equals(this.fieldMapOfStringToString, o.fieldMapOfStringToString)) return false;
        if (!Objects.equals(this.fieldAnonymousStruct, o.fieldAnonymousStruct)) return false;
        if (!Objects.equals(this.fieldRefToConstant, o.fieldRefToConstant)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldRef, this.fieldDisjunctionOfScalars, this.fieldMixedDisjunction, this.fieldDisjunctionWithNull, this.operator, this.fieldArrayOfStrings, this.fieldMapOfStringToString, this.fieldAnonymousStruct, this.fieldRefToConstant);
    }
}
