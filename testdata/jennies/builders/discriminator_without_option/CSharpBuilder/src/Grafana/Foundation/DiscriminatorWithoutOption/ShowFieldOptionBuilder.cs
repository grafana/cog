namespace Grafana.Foundation.DiscriminatorWithoutOption;


public class ShowFieldOptionBuilder : Cog.IBuilder<ShowFieldOption>
{
    protected readonly ShowFieldOption @internal;

    public ShowFieldOptionBuilder()
    {
        this.@internal = new ShowFieldOption();
    }

    public ShowFieldOptionBuilder Field(AnEnum field)
    {
        this.@internal.Field = field;
        return this;
    }

    public ShowFieldOptionBuilder Text(string text)
    {
        this.@internal.Text = text;
        return this;
    }

    public ShowFieldOption Build()
    {
        return this.@internal;
    }
}
