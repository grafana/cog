package struct_complex_fields;

import java.util.List;

// This struct does things.
public class SomeStruct {
    public SomeOtherStruct FieldRef;
    public StringOrBool FieldDisjunctionOfScalars;
    public StringOrSomeOtherStruct FieldMixedDisjunction;
    public StringOrNull FieldDisjunctionWithNull;
    public SomeStructOperator Operator;
    public List<String> FieldArrayOfStrings;
    public Map<String, String> FieldMapOfStringToString;
    public StructComplexFieldsSomeStructFieldAnonymousStruct FieldAnonymousStruct;
    
}
