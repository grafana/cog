package variant_panelcfg_full;

import java.util.Objects;

public class Options {
    public String timeseriesOption;
    public Options() {
        this.timeseriesOption = "";
    }
    public Options(String timeseriesOption) {
        this.timeseriesOption = timeseriesOption;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Options)) return false;
        Options o = (Options) other;
        if (!Objects.equals(this.timeseriesOption, o.timeseriesOption)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.timeseriesOption);
    }
}
