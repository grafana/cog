package enums;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum LogsSortOrder {
    ASC("time_asc"),
    DESC("time_desc"),
    _EMPTY("");

    private final String value;

    private LogsSortOrder(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
