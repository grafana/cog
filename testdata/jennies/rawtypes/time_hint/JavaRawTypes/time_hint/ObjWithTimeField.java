package time_hint;

import java.util.Objects;

public class ObjWithTimeField {
    public String registeredAt;
    public String duration;
    public ObjWithTimeField() {
        this.registeredAt = "";
        this.duration = "";
    }
    public ObjWithTimeField(String registeredAt,String duration) {
        this.registeredAt = registeredAt;
        this.duration = duration;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof ObjWithTimeField)) return false;
        ObjWithTimeField o = (ObjWithTimeField) other;
        if (!Objects.equals(this.registeredAt, o.registeredAt)) return false;
        if (!Objects.equals(this.duration, o.duration)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.registeredAt, this.duration);
    }
}
