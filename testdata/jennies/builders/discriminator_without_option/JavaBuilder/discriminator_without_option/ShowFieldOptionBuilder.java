package discriminator_without_option;


public class ShowFieldOptionBuilder implements cog.Builder<ShowFieldOption> {
    protected final ShowFieldOption internal;
    
    public ShowFieldOptionBuilder() {
        this.internal = new ShowFieldOption();
    }
    public ShowFieldOptionBuilder field(AnEnum field) {
        this.internal.field = field;
        return this;
    }
    
    public ShowFieldOptionBuilder text(String text) {
        this.internal.text = text;
        return this;
    }
    public ShowFieldOption build() {
        return this.internal;
    }
}
