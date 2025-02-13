import * as cog from '../cog';
import * as somePkg from '../somePkg';

export class SomeNiceBuilderBuilder implements cog.Builder<somePkg.SomeStruct> {
    protected readonly internal: somePkg.SomeStruct;

    constructor() {
        this.internal = somePkg.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): somePkg.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}

