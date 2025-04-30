package variant_dataquery;


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
}
