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
        if (!(id >= 5)) {
            throw new IllegalArgumentException("id must be >= 5");
        }
        if (!(id < 10)) {
            throw new IllegalArgumentException("id must be < 10");
        }
    this.internal.id = id;
        return this;
    }
    
    public Builder setTitle(String title) {
        if (!(title.length() >= 1)) {
            throw new IllegalArgumentException("title.length() must be >= 1");
        }
    this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
