namespace Grafana.Foundation.ComposableSlot;


public class LokiBuilderBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public LokiBuilderBuilder()
    {
        this.@internal = new Dashboard();
    }

    public LokiBuilderBuilder Target(Cog.IBuilder<Cog.Variants.Dataquery> target)
    {
        Cog.Variants.Dataquery targetResource = target.Build();
        this.@internal.Target = targetResource;
        return this;
    }

    public LokiBuilderBuilder Targets(List<Cog.IBuilder<Cog.Variants.Dataquery>> targets)
    {
        List<Cog.Variants.Dataquery> targetsResources = new List<Cog.Variants.Dataquery>();
        foreach (Cog.IBuilder<Cog.Variants.Dataquery> r1 in targets)
        {
                Cog.Variants.Dataquery targetsDepth1 = r1.Build();
                targetsResources.Add(targetsDepth1);
        }
        this.@internal.Targets = targetsResources;
        return this;
    }

    public Dashboard Build()
    {
        return this.@internal;
    }
}
