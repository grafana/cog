package discriminator_without_option;


public class NoShowFieldOptionBuilder implements cog.Builder<NoShowFieldOption> {
    protected final NoShowFieldOption internal;
    
    public NoShowFieldOptionBuilder() {
        this.internal = new NoShowFieldOption();
    }
    public NoShowFieldOptionBuilder text(String text) {
        this.internal.text = text;
        return this;
    }
    public NoShowFieldOption build() {
        return this.internal;
    }
}
