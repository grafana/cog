namespace Grafana.Foundation.DataqueryVariantBuilder;


public class LokiBuilderBuilder : Cog.IBuilder<Cog.Variants.Dataquery>
{
    protected readonly Loki @internal;

    public LokiBuilderBuilder()
    {
        this.@internal = new Loki();
    }

    public LokiBuilderBuilder Expr(string expr)
    {
        this.@internal.Expr = expr;
        return this;
    }

    public Loki Build()
    {
        return this.@internal;
    }
}
