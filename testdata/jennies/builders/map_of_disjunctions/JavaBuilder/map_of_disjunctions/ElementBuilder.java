package map_of_disjunctions;


public class ElementBuilder implements cog.Builder<Element> {
    protected final Element internal;
    
    public ElementBuilder() {
        this.internal = new Element();
    }
    public ElementBuilder panel(cog.Builder<Panel> panel) {
    Panel panelResource = panel.build();
        this.internal.panel = panelResource;
        return this;
    }
    
    public ElementBuilder libraryPanel(cog.Builder<LibraryPanel> libraryPanel) {
    LibraryPanel libraryPanelResource = libraryPanel.build();
        this.internal.libraryPanel = libraryPanelResource;
        return this;
    }
    public Element build() {
        return this.internal;
    }
}
