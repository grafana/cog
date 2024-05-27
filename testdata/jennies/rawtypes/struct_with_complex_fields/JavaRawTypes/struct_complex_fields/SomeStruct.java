package struct_complex_fields;

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
}
