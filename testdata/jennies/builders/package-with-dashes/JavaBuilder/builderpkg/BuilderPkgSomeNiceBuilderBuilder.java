package builder-pkg;

import with-dashes.SomeStruct;

public class BuilderPkgSomeNiceBuilderBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public BuilderPkgSomeNiceBuilderBuilder() {
        this.internal = new SomeStruct();
    }
    public BuilderPkgSomeNiceBuilderBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
