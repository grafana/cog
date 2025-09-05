package intersections;


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
}
