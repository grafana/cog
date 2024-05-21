package initialization_safeguards;


public class SomePanel {
    public String title;
    public Options options;
    
    public static class Builder {
        private SomePanel internal;
        
        public Builder() {
            this.internal = new SomePanel();
        }
    public Builder setTitle(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder setShowLegend(Boolean show) {
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
}
