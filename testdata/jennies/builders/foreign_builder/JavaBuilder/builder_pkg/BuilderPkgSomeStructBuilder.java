package builder_pkg;

import some_pkg.SomeStruct;

public class BuilderPkgSomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public BuilderPkgSomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public BuilderPkgSomeStructBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
