namespace Grafana.Foundation.Promql;


public class FuncCallExprBuilder : Cog.IBuilder<FuncCallExpr>
{
    protected readonly FuncCallExpr @internal;

    public FuncCallExprBuilder()
    {
        this.@internal = new FuncCallExpr();
        this.@internal.Type = "funcCallExpr";
    }

    public FuncCallExprBuilder Function(string function)
    {
        if (!(function.Length >= 1))
        {
            throw new System.ArgumentException("function.Length must be >= 1");
        }
        this.@internal.Function = function;
        return this;
    }

    public FuncCallExprBuilder Args(List<Cog.IBuilder<Expr>> args)
    {
        List<Expr> argsResources = new List<Expr>();
        foreach (Cog.IBuilder<Expr> r1 in args)
        {
                Expr argsDepth1 = r1.Build();
                argsResources.Add(argsDepth1);
        }
        this.@internal.Args = argsResources;
        return this;
    }

    public FuncCallExprBuilder Arg(Cog.IBuilder<Expr> arg)
    {
        if (this.@internal.Args == null)
        {
            this.@internal.Args = new List<Expr>();
        }
        Expr argResource = arg.Build();
        this.@internal.Args.Add(argResource);
        return this;
    }

    public FuncCallExpr Build()
    {
        return this.@internal;
    }
}
