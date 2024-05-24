package com.grafana.foundation.enums;


public enum TableSortOrder {
    ASC("asc"),
    DESC("desc");

    private String value;

    private TableSortOrder(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
