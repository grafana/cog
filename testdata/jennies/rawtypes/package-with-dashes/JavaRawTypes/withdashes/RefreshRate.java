package withdashes;


// Refresh rate or disabled.
public class RefreshRate {
    protected String string;
    protected Boolean bool;
    protected RefreshRate() {}
    public static RefreshRate createString(String string) {
        RefreshRate refreshRate = new RefreshRate();
        refreshRate.string = string;
        return refreshRate;
    }
    public static RefreshRate createBool(Boolean bool) {
        RefreshRate refreshRate = new RefreshRate();
        refreshRate.bool = bool;
        return refreshRate;
    }
}
