package dashboard;

import java.util.Objects;

public class FieldConfigSource {
    public FieldConfig defaults;
    public FieldConfigSource() {
    }
    public FieldConfigSource(FieldConfig defaults) {
        this.defaults = defaults;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof FieldConfigSource)) return false;
        FieldConfigSource o = (FieldConfigSource) other;
        if (!Objects.equals(this.defaults, o.defaults)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.defaults);
    }
}
