package builder-pkg;

import with-dashes.SomeStruct;

public class BuilderPkgSomeStructSomeNiceBuilderBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public BuilderPkgSomeStructSomeNiceBuilderBuilder() {
        this.internal = new SomeStruct();
    }
    public BuilderPkgSomeStructSomeNiceBuilderBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
