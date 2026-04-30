package variant_panelcfg_full;

import java.util.Objects;

public class FieldConfig {
    public String timeseriesFieldConfigOption;
    public FieldConfig() {
        this.timeseriesFieldConfigOption = "";
    }
    public FieldConfig(String timeseriesFieldConfigOption) {
        this.timeseriesFieldConfigOption = timeseriesFieldConfigOption;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof FieldConfig)) return false;
        FieldConfig o = (FieldConfig) other;
        if (!Objects.equals(this.timeseriesFieldConfigOption, o.timeseriesFieldConfigOption)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.timeseriesFieldConfigOption);
    }
}
