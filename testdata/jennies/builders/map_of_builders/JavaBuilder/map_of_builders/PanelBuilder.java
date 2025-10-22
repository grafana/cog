package map_of_builders;


public class PanelBuilder implements cog.Builder<Panel> {
    protected final Panel internal;
    
    public PanelBuilder() {
        this.internal = new Panel();
    }
    public PanelBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public Panel build() {
        return this.internal;
    }
}
