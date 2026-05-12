package dashboard;


public class PanelBuilder implements cog.Builder<Panel> {
    protected final Panel internal;
    
    public PanelBuilder() {
        this.internal = new Panel();
    }
    public PanelBuilder onlyFromThisDashboard(Boolean onlyFromThisDashboard) {
        this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }
    public Panel build() {
        return this.internal;
    }
}
