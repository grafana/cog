package sandbox;


public class SomeStruct {
    public unknown editable;
    public unknown autoRefresh;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Editable() {
    this.internal.editable = true;
        return this;
    }
    
    public Builder Readonly() {
    this.internal.editable = false;
        return this;
    }
    
    public Builder AutoRefresh() {
    this.internal.autoRefresh = true;
        return this;
    }
    
    public Builder NoAutoRefresh() {
    this.internal.autoRefresh = false;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
