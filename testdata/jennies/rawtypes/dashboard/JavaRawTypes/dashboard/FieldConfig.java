package dashboard;

import java.util.Objects;

public class FieldConfig {
    public String unit;
    public Object custom;
    public FieldConfig() {
    }
    public FieldConfig(String unit,Object custom) {
        this.unit = unit;
        this.custom = custom;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof FieldConfig)) return false;
        FieldConfig o = (FieldConfig) other;
        if (!Objects.equals(this.unit, o.unit)) return false;
        if (!Objects.equals(this.custom, o.custom)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.unit, this.custom);
    }
}
