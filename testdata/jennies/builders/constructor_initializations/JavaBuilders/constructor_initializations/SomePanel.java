package constructor_initializations;


public class SomePanel {
    public String type;
    public String title;
    public CursorMode cursor;
        
    public static class Builder {
        private SomePanel internal;
        
        public Builder() {
            this.internal = new SomePanel();
    this.internal.type = "panel_type";
    this.internal.cursor = CursorMode.TOOLTIP;
        }
    public Builder setTitle(String title) {
    this.internal.title = title;
        return this;
    }
    public SomePanel build() {
            return this.internal;
        }
    }
}
