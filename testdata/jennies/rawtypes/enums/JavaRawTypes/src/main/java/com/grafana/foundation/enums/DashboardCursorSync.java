package com.grafana.foundation.enums;


// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
public enum DashboardCursorSync {
    OFF(0),
    CROSSHAIR(1),
    TOOLTIP(2);

    private Integer value;

    private DashboardCursorSync(Integer value) {
        this.value = value;
    }

    public Integer getValue() {
        return value;
    }
}
