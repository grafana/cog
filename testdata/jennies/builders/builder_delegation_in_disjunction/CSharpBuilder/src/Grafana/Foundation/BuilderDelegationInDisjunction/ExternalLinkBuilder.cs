namespace Grafana.Foundation.BuilderDelegationInDisjunction;


public class ExternalLinkBuilder : Cog.IBuilder<ExternalLink>
{
    protected readonly ExternalLink @internal;

    public ExternalLinkBuilder()
    {
        this.@internal = new ExternalLink();
    }

    public ExternalLinkBuilder Url(string url)
    {
        this.@internal.Url = url;
        return this;
    }

    public ExternalLink Build()
    {
        return this.@internal;
    }
}
