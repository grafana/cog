import * as cog from '../cog';
import * as withDashes from '../withDashes';

export class SomeNiceBuilderBuilder implements cog.Builder<withDashes.SomeStruct> {
    protected readonly internal: withDashes.SomeStruct;

    constructor() {
        this.internal = withDashes.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): withDashes.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}

