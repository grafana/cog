import * as cog from '../cog';
import * as constraints from '../constraints';

export class SomeStructBuilder implements cog.Builder<constraints.SomeStruct> {
    protected readonly internal: constraints.SomeStruct;

    constructor() {
        this.internal = constraints.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): constraints.SomeStruct {
        return this.internal;
    }

    id(id: number): this {
        if (!(id >= 5)) {
            throw new Error("id must be >= 5");
        }
        if (!(id < 10)) {
            throw new Error("id must be < 10");
        }
        this.internal.id = id;
        return this;
    }

    title(title: string): this {
        if (!(title.length >= 1)) {
            throw new Error("title.length must be >= 1");
        }
        this.internal.title = title;
        return this;
    }
}

