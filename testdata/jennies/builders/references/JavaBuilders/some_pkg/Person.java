package some_pkg;

import other_pkg.Name;

public class Person {
    public Name name;
        
    public static class Builder {
        private Person internal;
        
        public Builder() {
            this.internal = new Person();
        }
    public Builder setName(Name name) {
    this.internal.name = name;
        return this;
    }
    public Person build() {
            return this.internal;
        }
    }
}
