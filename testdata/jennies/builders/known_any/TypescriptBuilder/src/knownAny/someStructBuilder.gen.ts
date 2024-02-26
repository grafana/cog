import * as cog from '../cog';
import * as knownAny from '../knownAny';

export class SomeStructBuilder implements cog.Builder<knownAny.SomeStruct> {
    private readonly internal: knownAny.SomeStruct;

    constructor() {
        this.internal = knownAny.defaultSomeStruct();
    }

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
