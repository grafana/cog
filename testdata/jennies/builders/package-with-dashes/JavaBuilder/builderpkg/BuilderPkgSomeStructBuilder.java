package builderpkg;

import withdashes.SomeStruct;

public class BuilderpkgSomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public BuilderpkgSomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public BuilderpkgSomeStructBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
