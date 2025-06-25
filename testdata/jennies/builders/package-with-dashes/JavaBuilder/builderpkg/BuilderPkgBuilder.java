package builder-pkg;

import with-dashes.SomeStruct;

public class BuilderPkgBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public BuilderPkgBuilder() {
        this.internal = new SomeStruct();
    }
    public BuilderPkgBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
