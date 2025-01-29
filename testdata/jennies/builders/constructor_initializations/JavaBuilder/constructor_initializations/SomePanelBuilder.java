package constructor_initializations;


public class SomePanelBuilder implements cog.Builder<SomePanel> {
    protected final SomePanel internal;
    
    public SomePanelBuilder() {
        this.internal = new SomePanel();
        this.internal.type = "panel_type";
        this.internal.cursor = CursorMode.TOOLTIP;
    }
    public SomePanelBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomePanel build() {
        return this.internal;
    }
}
