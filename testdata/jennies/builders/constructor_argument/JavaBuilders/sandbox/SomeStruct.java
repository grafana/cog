package sandbox;


public class SomeStruct {
    public String title;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder(String title) {
            this.internal = new SomeStruct();
    this.internal.title = title;
        }
    public Builder Title(String title) {
    this.internal.title = title;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
