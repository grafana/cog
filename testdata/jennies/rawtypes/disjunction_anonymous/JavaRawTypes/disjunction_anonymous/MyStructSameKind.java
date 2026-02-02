package disjunction_anonymous;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonValue;


@JsonFormat(shape = JsonFormat.Shape.OBJECT)
public enum MyStructSameKind {
    A("a"),
    B("b"),
    C("c"),
    _EMPTY("");

    private final String value;

    private MyStructSameKind(String value) {
        this.value = value;
    }

    @JsonValue
    public String Value() {
        return value;
    }
}
