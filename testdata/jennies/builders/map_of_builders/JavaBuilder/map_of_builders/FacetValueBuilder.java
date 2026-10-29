package map_of_builders;


public class FacetValueBuilder implements cog.Builder<FacetValue> {
    protected final FacetValue internal;
    
    public FacetValueBuilder() {
        this.internal = new FacetValue();
    }
    public FacetValueBuilder value(String value) {
        this.internal.value = value;
        return this;
    }
    
    public FacetValueBuilder count(Long count) {
        this.internal.count = count;
        return this;
    }
    public FacetValue build() {
        return this.internal;
    }
}
