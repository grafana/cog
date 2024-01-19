import * as cog from '../cog';
import * as known_any from '../known_any';

export class SomeStructBuilder implements cog.Builder<known_any.SomeStruct> {
    private readonly internal: known_any.SomeStruct;

    constructor() {
        this.internal = known_any.defaultSomeStruct();
    }

    build(): known_any.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        if (!this.internal.config) {
            this.internal.config = known_any.defaultConfig();
        }
        this.internal.config.title = title;
        return this;
    }
}
