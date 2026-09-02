namespace Grafana.Foundation.Sandbox;


public class DashboardBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public DashboardBuilder()
    {
        this.@internal = new Dashboard();
    }

    public DashboardBuilder WithVariable(string name,string value)
    {
        if (this.@internal.Variables == null)
        {
            this.@internal.Variables = new List<Variable>();
        }
    Variable variable = new Variable();
        variable.Name = name;
        variable.Value = value;
        this.@internal.Variables.Add(variable);
        return this;
    }

    public Dashboard Build()
    {
        return this.@internal;
    }
}
