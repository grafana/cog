import * as cog from '../cog';
import * as some_pkg from '../some_pkg';
import * as other_pkg from '../other_pkg';

export class PersonBuilder implements cog.Builder<some_pkg.Person> {
    private readonly internal: some_pkg.Person;

    constructor() {
        this.internal = some_pkg.defaultPerson();
    }

    build(): some_pkg.Person {
        return this.internal;
    }

    name(name: other_pkg.Name): this {
        this.internal.name = name;
        return this;
    }
}
