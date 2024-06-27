package some_pkg;

import other_pkg.Name;

public class Person {
    public Name name;
    
    public static class Builder implements cog.Builder<Person> {
        private Person internal;
        
        public Builder() {
            this.internal = new Person();
        }
    public Builder Name(Name name) {
    this.internal.name = name;
        return this;
    }
    public Person Build() {
            return this.internal;
        }
    }
}
