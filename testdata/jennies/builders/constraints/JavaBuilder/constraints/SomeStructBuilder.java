package constraints;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder id(Long id) {
        if (!(id >= 5)) {
            throw new IllegalArgumentException("id must be >= 5");
        }
        if (!(id < 10)) {
            throw new IllegalArgumentException("id must be < 10");
        }
        this.internal.id = id;
        return this;
    }
    
    public SomeStructBuilder title(String title) {
        if (!(title.length() >= 1)) {
            throw new IllegalArgumentException("title.length() must be >= 1");
        }
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
