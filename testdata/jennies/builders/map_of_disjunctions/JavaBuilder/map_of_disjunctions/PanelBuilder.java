package map_of_disjunctions;


public class PanelBuilder<T extends PanelBuilder<T>> implements cog.Builder<Panel> {
    protected final Panel internal;
    
    public PanelBuilder() {
        this.internal = new Panel();
        this.internal.kind = "Panel";
    }
    public T title(String title) {
        this.internal.title = title;
        return (T) this;
    }
    public Panel build() {
        return this.internal;
    }
}
