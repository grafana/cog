package disjunctions;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof RefreshRate)) return false;
        RefreshRate o = (RefreshRate) other;
        if (!Objects.equals(this.string, o.string)) return false;
        if (!Objects.equals(this.bool, o.bool)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.string, this.bool);
    }
}
