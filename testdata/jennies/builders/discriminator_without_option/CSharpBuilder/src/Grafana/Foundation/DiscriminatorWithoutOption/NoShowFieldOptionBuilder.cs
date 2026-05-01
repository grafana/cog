namespace Grafana.Foundation.DiscriminatorWithoutOption;


public class NoShowFieldOptionBuilder : Cog.IBuilder<NoShowFieldOption>
{
    protected readonly NoShowFieldOption @internal;

    public NoShowFieldOptionBuilder()
    {
        this.@internal = new NoShowFieldOption();
    }

    public NoShowFieldOptionBuilder Text(string text)
    {
        this.@internal.Text = text;
        return this;
    }

    public NoShowFieldOption Build()
    {
        return this.@internal;
    }
}
