package intersections;

import java.util.Objects;

// Base properties for all metrics
public class Common {
    // The metric name
    public String name;
    // The metric type
    public CommonType type;
    // The type of data the metric contains
    public CommonContains contains;
    public Common() {
        this.name = "";
        this.type = CommonType.COUNTER;
        this.contains = CommonContains.DEFAULT;
    }
    public Common(String name,CommonType type,CommonContains contains) {
        this.name = name;
        this.type = type;
        this.contains = contains;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Common)) return false;
        Common o = (Common) other;
        if (!Objects.equals(this.name, o.name)) return false;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.contains, o.contains)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.name, this.type, this.contains);
    }
}
