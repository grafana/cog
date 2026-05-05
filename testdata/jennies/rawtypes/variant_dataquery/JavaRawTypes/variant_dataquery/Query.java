package variant_dataquery;

import java.util.Objects;

public class Query implements cog.variants.Dataquery {
    public String expr;
    public Boolean instant;
    public Query() {
        this.expr = "";
    }
    public Query(String expr,Boolean instant) {
        this.expr = expr;
        this.instant = instant;
    }
    public String dataqueryName() {
        return "prometheus";
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Query)) return false;
        Query o = (Query) other;
        if (!Objects.equals(this.expr, o.expr)) return false;
        if (!Objects.equals(this.instant, o.instant)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.expr, this.instant);
    }
}
