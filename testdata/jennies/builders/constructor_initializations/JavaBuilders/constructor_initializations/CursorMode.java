package constructor_initializations;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum CursorMode {
    OFF("off"),
    TOOLTIP("tooltip"),
    CROSSHAIR("crosshair");

    private final String value;

    private CursorMode(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
