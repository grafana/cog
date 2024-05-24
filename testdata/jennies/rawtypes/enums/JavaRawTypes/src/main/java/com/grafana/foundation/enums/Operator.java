package com.grafana.foundation.enums;


// This is a very interesting string enum.
public enum Operator {
    GREATER_THAN(">"),
    LESS_THAN("<");

    private String value;

    private Operator(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
