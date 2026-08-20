namespace Grafana.Foundation.VariantDataquery;


public class Query
{
    public string Expr;
    public bool Instant;

    public Query()
    {
        this.Expr = "";
    }

    public Query(string expr, bool instant)
    {
        this.Expr = expr;
        this.Instant = instant;
    }
}
