package com.grafana.foundation.constructor_initializations;


public enum CursorMode {
    OFF("off"),
    TOOLTIP("tooltip"),
    CROSSHAIR("crosshair");

    private String value;

    private CursorMode(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
