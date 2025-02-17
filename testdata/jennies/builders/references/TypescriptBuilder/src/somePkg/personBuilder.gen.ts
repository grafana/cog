import * as cog from '../cog';
import * as somePkg from '../somePkg';
import * as otherPkg from '../otherPkg';

export class PersonBuilder implements cog.Builder<somePkg.Person> {
    protected readonly internal: somePkg.Person;

    constructor() {
        this.internal = somePkg.defaultPerson();
    }

    /**
     * Builds the object.
     */
    build(): somePkg.Person {
        return this.internal;
    }

    name(name: otherPkg.Name): this {
        this.internal.name = name;
        return this;
    }
}

