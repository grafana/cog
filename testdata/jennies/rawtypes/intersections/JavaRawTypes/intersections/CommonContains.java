package intersections;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum CommonContains {
    DEFAULT("default"),
    TIME("time"),
    _EMPTY("");

    private final String value;

    private CommonContains(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
