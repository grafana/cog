namespace Grafana.Foundation.MapOfDisjunctions;


public class LibraryPanelBuilder : Cog.IBuilder<LibraryPanel>
{
    protected readonly LibraryPanel @internal;

    public LibraryPanelBuilder()
    {
        this.@internal = new LibraryPanel();
        this.@internal.Kind = "Library";
    }

    public LibraryPanelBuilder Text(string text)
    {
        this.@internal.Text = text;
        return this;
    }

    public LibraryPanel Build()
    {
        return this.@internal;
    }
}
