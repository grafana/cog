namespace Grafana.Foundation.MapOfDisjunctions;


public class ElementBuilder : Cog.IBuilder<Element>
{
    protected readonly Element @internal;

    public ElementBuilder()
    {
        this.@internal = new Element();
    }

    public ElementBuilder Panel(Cog.IBuilder<Panel> panel)
    {
        Panel panelResource = panel.Build();
        this.@internal.Panel = panelResource;
        return this;
    }

    public ElementBuilder LibraryPanel(Cog.IBuilder<LibraryPanel> libraryPanel)
    {
        LibraryPanel libraryPanelResource = libraryPanel.Build();
        this.@internal.LibraryPanel = libraryPanelResource;
        return this;
    }

    public Element Build()
    {
        return this.@internal;
    }
}
