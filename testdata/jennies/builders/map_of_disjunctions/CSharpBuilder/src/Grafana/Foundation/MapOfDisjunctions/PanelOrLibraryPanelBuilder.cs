namespace Grafana.Foundation.MapOfDisjunctions;


public class PanelOrLibraryPanelBuilder : Cog.IBuilder<PanelOrLibraryPanel>
{
    protected readonly PanelOrLibraryPanel @internal;

    public PanelOrLibraryPanelBuilder()
    {
        this.@internal = new PanelOrLibraryPanel();
    }

    public PanelOrLibraryPanelBuilder Panel(Cog.IBuilder<Panel> panel)
    {
        Panel panelResource = panel.Build();
        this.@internal.Panel = panelResource;
        return this;
    }

    public PanelOrLibraryPanelBuilder LibraryPanel(Cog.IBuilder<LibraryPanel> libraryPanel)
    {
        LibraryPanel libraryPanelResource = libraryPanel.Build();
        this.@internal.LibraryPanel = libraryPanelResource;
        return this;
    }

    public PanelOrLibraryPanel Build()
    {
        return this.@internal;
    }
}
