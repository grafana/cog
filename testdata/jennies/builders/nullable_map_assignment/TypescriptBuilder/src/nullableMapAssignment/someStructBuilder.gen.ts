import * as cog from '../cog';
import * as nullableMapAssignment from '../nullableMapAssignment';

export class SomeStructBuilder implements cog.Builder<nullableMapAssignment.SomeStruct> {
    protected readonly internal: nullableMapAssignment.SomeStruct;

    constructor() {
        this.internal = nullableMapAssignment.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): nullableMapAssignment.SomeStruct {
        return this.internal;
    }

    config(config: Record<string, string>): this {
        this.internal.config = config;
        return this;
    }
}

