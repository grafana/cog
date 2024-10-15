package variant_dataquery;


public class Query implements cog.variants.Dataquery {
    public String expr;
    public Boolean instant;
    public String dataqueryName() {
        return "prometheus";
    }
}
