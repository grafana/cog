package map_of_disjunctions;


public class LibraryPanelBuilder implements cog.Builder<LibraryPanel> {
    protected final LibraryPanel internal;
    
    public LibraryPanelBuilder() {
        this.internal = new LibraryPanel();
        this.internal.kind = "Library";
    }
    public LibraryPanelBuilder text(String text) {
        this.internal.text = text;
        return this;
    }
    public LibraryPanel build() {
        return this.internal;
    }
}
