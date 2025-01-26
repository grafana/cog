package initialization_safeguards;


public class SomePanelBuilder implements cog.Builder<SomePanel> {
    protected final SomePanel internal;
    
    public SomePanelBuilder() {
        this.internal = new SomePanel();
    }
    public SomePanelBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    
    public SomePanelBuilder showLegend(Boolean show) {
		if (this.internal.options == null) {
			this.internal.options = new initialization_safeguards.Options();
		}
		if (this.internal.options.legend == null) {
			this.internal.options.legend = new initialization_safeguards.LegendOptions();
		}
        this.internal.options.legend.show = show;
        return this;
    }
    public SomePanel build() {
        return this.internal;
    }
}
