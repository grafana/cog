namespace Grafana.Foundation.StructWithDefaults;


public class StructBuilder : Cog.IBuilder<Struct>
{
    protected readonly Struct @internal;

    public StructBuilder()
    {
        this.@internal = new Struct();
    }

    public StructBuilder AllFields(Cog.IBuilder<NestedStruct> allFields)
    {
        NestedStruct allFieldsResource = allFields.Build();
        this.@internal.AllFields = allFieldsResource;
        return this;
    }

    public StructBuilder PartialFields(Cog.IBuilder<NestedStruct> partialFields)
    {
        NestedStruct partialFieldsResource = partialFields.Build();
        this.@internal.PartialFields = partialFieldsResource;
        return this;
    }

    public StructBuilder EmptyFields(Cog.IBuilder<NestedStruct> emptyFields)
    {
        NestedStruct emptyFieldsResource = emptyFields.Build();
        this.@internal.EmptyFields = emptyFieldsResource;
        return this;
    }

    public StructBuilder ComplexField(object complexField)
    {
        this.@internal.ComplexField = complexField;
        return this;
    }

    public StructBuilder PartialComplexField(object partialComplexField)
    {
        this.@internal.PartialComplexField = partialComplexField;
        return this;
    }

    public Struct Build()
    {
        return this.@internal;
    }
}
