package dashboard;


public class PanelBuilder<T extends PanelBuilder<T>> implements cog.Builder<Panel> {
    protected final Panel internal;
    
    public PanelBuilder() {
        this.internal = new Panel();
    }
    public T onlyFromThisDashboard(Boolean onlyFromThisDashboard) {
        this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return (T) this;
    }
    public Panel build() {
        return this.internal;
    }
}
