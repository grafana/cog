package enums;


public enum LogsSortOrder {
    ASC("time_asc"),
    DESC("time_desc");

    private String value;

    private LogsSortOrder(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
