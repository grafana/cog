package struct_optional_fields;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum SomeStructOperator {
    GREATER_THAN(">"),
    LESS_THAN("<"),
    _EMPTY("");

    private final String value;

    private SomeStructOperator(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
