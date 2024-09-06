package struct_with_defaults;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonInclude;

public class Struct {
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("allFields")
    public NestedStruct allFields;
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("partialFields")
    public NestedStruct partialFields;
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("emptyFields")
    public NestedStruct emptyFields;
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("complexField")
    public Object complexField;
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("partialComplexField")
    public Object partialComplexField;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Struct> {
        private final Struct internal;
        
        public Builder() {
            this.internal = new Struct();
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.stringVal("hello");
        nestedStructResource.intVal(3L);
        this.allFields(nestedStructResource);
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.intVal(4L);
        this.partialFields(nestedStructResource);
        this.complexField(new Object());
        this.partialComplexField(new Object());
        }
    public Builder allFields(cog.Builder<NestedStruct> allFields) {
    this.internal.allFields = allFields.build();
        return this;
    }
    
    public Builder partialFields(cog.Builder<NestedStruct> partialFields) {
    this.internal.partialFields = partialFields.build();
        return this;
    }
    
    public Builder emptyFields(cog.Builder<NestedStruct> emptyFields) {
    this.internal.emptyFields = emptyFields.build();
        return this;
    }
    
    public Builder complexField(Object complexField) {
    this.internal.complexField = complexField;
        return this;
    }
    
    public Builder partialComplexField(Object partialComplexField) {
    this.internal.partialComplexField = partialComplexField;
        return this;
    }
    public Struct build() {
            return this.internal;
        }
    }
}
