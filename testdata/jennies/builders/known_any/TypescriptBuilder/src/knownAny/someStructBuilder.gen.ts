import * as cog from '../cog';
import * as knownAny from '../knownAny';

export class SomeStructBuilder implements cog.Builder<knownAny.SomeStruct> {
    protected readonly internal: knownAny.SomeStruct;

    constructor() {
        this.internal = knownAny.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): knownAny.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        if (!this.internal.config) {
            this.internal.config = knownAny.defaultConfig();
        }
        this.internal.config.title = title;
        return this;
    }
}

