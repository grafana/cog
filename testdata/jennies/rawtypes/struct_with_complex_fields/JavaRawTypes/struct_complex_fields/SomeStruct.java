package struct_complex_fields;

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
}
