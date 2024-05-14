package constraints;


public class SomeStruct {
    public Long id;
    public String title;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setId(Long id) {
        if (id >= 5) {
            return this;
        }
        if (id < 10) {
            return this;
        }
    this.internal.id = id;
        return this;
    }
    
    public Builder setTitle(String title) {
        if (title.length() >= 1) {
            return this;
        }
    this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
