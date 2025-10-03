package map_of_disjunctions;


public class PanelOrLibraryPanelBuilder implements cog.Builder<PanelOrLibraryPanel> {
    protected final PanelOrLibraryPanel internal;
    
    public PanelOrLibraryPanelBuilder() {
        this.internal = new PanelOrLibraryPanel();
    }
    public PanelOrLibraryPanelBuilder panel(cog.Builder<Panel> panel) {
    Panel panelResource = panel.build();
        this.internal.panel = panelResource;
        return this;
    }
    
    public PanelOrLibraryPanelBuilder libraryPanel(cog.Builder<LibraryPanel> libraryPanel) {
    LibraryPanel libraryPanelResource = libraryPanel.build();
        this.internal.libraryPanel = libraryPanelResource;
        return this;
    }
    public PanelOrLibraryPanel build() {
        return this.internal;
    }
}
