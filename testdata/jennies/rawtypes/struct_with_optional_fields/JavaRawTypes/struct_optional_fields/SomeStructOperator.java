package struct_optional_fields;


public enum SomeStructOperator {
    GREATER_THAN(">"),
    LESS_THAN("<");

    private String value;

    private SomeStructOperator(String value) {
        this.value = value;
    }

    public String Value() {
        return value;
    }
}
