namespace Grafana.Foundation.Enums;

// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
public enum DashboardCursorSync
{
    Off = 0,
    Crosshair = 1,
    Tooltip = 2
}
