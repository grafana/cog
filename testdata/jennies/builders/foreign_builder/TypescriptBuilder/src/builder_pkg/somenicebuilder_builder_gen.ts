import * as cog from '../cog';
import * as some_pkg from '../some_pkg';

export class SomeNiceBuilderBuilder implements cog.Builder<some_pkg.SomeStruct> {
    private readonly internal: some_pkg.SomeStruct;

    constructor() {
        this.internal = some_pkg.defaultSomeStruct();
    }

    build(): some_pkg.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
