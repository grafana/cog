package some_pkg;

import other_pkg.Name;

public class PersonBuilder implements cog.Builder<Person> {
    protected final Person internal;
    
    public PersonBuilder() {
        this.internal = new Person();
    }
    public PersonBuilder name(Name name) {
        this.internal.name = name;
        return this;
    }
    public Person build() {
        return this.internal;
    }
}
