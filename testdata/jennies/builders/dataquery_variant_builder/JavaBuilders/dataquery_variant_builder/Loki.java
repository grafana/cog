package dataquery_variant_builder;


public class Loki implements cog.variants.Dataquery {
    public String expr;
    
    public static class Builder implements cog.Builder<Loki> {
        private Loki internal;
        
        public Builder() {
            this.internal = new Loki();
        }
    public Builder Expr(String expr) {
    this.internal.expr = expr;
        return this;
    }
    public Loki build() {
            return this.internal;
        }
    }
}
