package dataquery_variant_builder;


public class LokiBuilderBuilder implements cog.Builder<cog.variants.Dataquery> {
    protected final Loki internal;
    
    public LokiBuilderBuilder() {
        this.internal = new Loki();
    }
    public LokiBuilderBuilder expr(String expr) {
        this.internal.expr = expr;
        return this;
    }
    public Loki build() {
        return this.internal;
    }
}
