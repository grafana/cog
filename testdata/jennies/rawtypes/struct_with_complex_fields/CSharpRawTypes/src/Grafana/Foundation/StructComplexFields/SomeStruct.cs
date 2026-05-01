namespace Grafana.Foundation.StructComplexFields;

using System.Collections.Generic;

// This struct does things.
public class SomeStruct
{
    public SomeOtherStruct FieldRef;
    public StringOrBool FieldDisjunctionOfScalars;
    public StringOrSomeOtherStruct FieldMixedDisjunction;
    public string FieldDisjunctionWithNull;
    public SomeStructOperator Operator;
    public List<string> FieldArrayOfStrings;
    public Dictionary<string, string> FieldMapOfStringToString;
    public StructComplexFieldsSomeStructFieldAnonymousStruct FieldAnonymousStruct;
    public string FieldRefToConstant;

    public SomeStruct()
    {
        this.FieldRef = new SomeOtherStruct();
        this.FieldDisjunctionOfScalars = new StringOrBool();
        this.FieldMixedDisjunction = new StringOrSomeOtherStruct();
        this.Operator = SomeStructOperator.GreaterThan;
        this.FieldArrayOfStrings = new List<string>();
        this.FieldMapOfStringToString = new Dictionary<string, string>();
        this.FieldAnonymousStruct = new StructComplexFieldsSomeStructFieldAnonymousStruct();
        this.FieldRefToConstant = new ConnectionPath();
    }

    public SomeStruct(SomeOtherStruct fieldRef, StringOrBool fieldDisjunctionOfScalars, StringOrSomeOtherStruct fieldMixedDisjunction, string fieldDisjunctionWithNull, SomeStructOperator operatorArg, List<string> fieldArrayOfStrings, Dictionary<string, string> fieldMapOfStringToString, StructComplexFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct, string fieldRefToConstant)
    {
        this.FieldRef = fieldRef;
        this.FieldDisjunctionOfScalars = fieldDisjunctionOfScalars;
        this.FieldMixedDisjunction = fieldMixedDisjunction;
        this.FieldDisjunctionWithNull = fieldDisjunctionWithNull;
        this.Operator = operatorArg;
        this.FieldArrayOfStrings = fieldArrayOfStrings;
        this.FieldMapOfStringToString = fieldMapOfStringToString;
        this.FieldAnonymousStruct = fieldAnonymousStruct;
        this.FieldRefToConstant = fieldRefToConstant;
    }
}
