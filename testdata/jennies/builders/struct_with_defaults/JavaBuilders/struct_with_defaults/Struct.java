package struct_with_defaults;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Struct { 
    @JsonProperty("allFields")
    public NestedStruct allFields; 
    @JsonProperty("partialFields")
    public NestedStruct partialFields; 
    @JsonProperty("emptyFields")
    public NestedStruct emptyFields; 
    @JsonProperty("complexField")
    public Object complexField; 
    @JsonProperty("partialComplexField")
    public Object partialComplexField;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Struct> {
        private Struct internal;
        
        public Builder() {
            this.internal = new Struct();
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.StringVal("hello");
        nestedStructResource.IntVal(3L);
        this.AllFields(nestedStructResource);
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.IntVal(4L);
        this.PartialFields(nestedStructResource);
        this.ComplexField(new Object());
        this.PartialComplexField(new Object());
        }
    public Builder AllFields(cog.Builder<NestedStruct> allFields) {
    this.internal.allFields = allFields.Build();
        return this;
    }
    
    public Builder PartialFields(cog.Builder<NestedStruct> partialFields) {
    this.internal.partialFields = partialFields.Build();
        return this;
    }
    
    public Builder EmptyFields(cog.Builder<NestedStruct> emptyFields) {
    this.internal.emptyFields = emptyFields.Build();
        return this;
    }
    
    public Builder ComplexField(Object complexField) {
    this.internal.complexField = complexField;
        return this;
    }
    
    public Builder PartialComplexField(Object partialComplexField) {
    this.internal.partialComplexField = partialComplexField;
        return this;
    }
    public Struct Build() {
            return this.internal;
        }
    }
}
