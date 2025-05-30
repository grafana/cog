package builder-pkg;

import with-dashes.SomeStruct;

public class SomeNiceBuilderBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeNiceBuilderBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeNiceBuilderBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
