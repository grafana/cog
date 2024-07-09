package enums;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


// This is a very interesting string enum.
@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum Operator {
    GREATER_THAN(">"),
    LESS_THAN("<"),
    _EMPTY("");

    private final String value;

    private Operator(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
