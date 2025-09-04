package intersections;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum CommonType {
    COUNTER("counter"),
    GAUGE("gauge"),
    _EMPTY("");

    private final String value;

    private CommonType(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
