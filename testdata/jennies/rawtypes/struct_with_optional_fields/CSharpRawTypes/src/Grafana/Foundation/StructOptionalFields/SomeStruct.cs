namespace Grafana.Foundation.StructOptionalFields;

using System.Collections.Generic;

public class SomeStruct
{
    public SomeOtherStruct FieldRef;
    public string FieldString;
    public SomeStructOperator Operator;
    public List<string> FieldArrayOfStrings;
    public StructOptionalFieldsSomeStructFieldAnonymousStruct FieldAnonymousStruct;

    public SomeStruct()
    {
    }

    public SomeStruct(SomeOtherStruct fieldRef, string fieldString, SomeStructOperator operatorArg, List<string> fieldArrayOfStrings, StructOptionalFieldsSomeStructFieldAnonymousStruct fieldAnonymousStruct)
    {
        this.FieldRef = fieldRef;
        this.FieldString = fieldString;
        this.Operator = operatorArg;
        this.FieldArrayOfStrings = fieldArrayOfStrings;
        this.FieldAnonymousStruct = fieldAnonymousStruct;
    }
}
