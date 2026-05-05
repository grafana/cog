package variant_panelcfg_only_options;

import java.util.Objects;

public class Options {
    public String content;
    public Options() {
        this.content = "";
    }
    public Options(String content) {
        this.content = content;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Options)) return false;
        Options o = (Options) other;
        if (!Objects.equals(this.content, o.content)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.content);
    }
}
