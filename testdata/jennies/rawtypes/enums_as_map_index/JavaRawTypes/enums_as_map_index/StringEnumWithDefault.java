package enums_as_map_index;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum StringEnumWithDefault {
    A("a"),
    B("b"),
    C("c"),
    _EMPTY("");

    private final String value;

    private StringEnumWithDefault(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
