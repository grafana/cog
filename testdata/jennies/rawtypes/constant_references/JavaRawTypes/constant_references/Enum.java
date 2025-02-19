package constant_references;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum Enum {
    VALUE_A("ValueA"),
    VALUE_B("ValueB"),
    VALUE_C("ValueC"),
    _EMPTY("");

    private final String value;

    private Enum(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
