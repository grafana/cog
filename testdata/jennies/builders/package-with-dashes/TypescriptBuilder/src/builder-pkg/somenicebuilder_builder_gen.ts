import * as cog from '../cog';
import * as with-dashes from '../with-dashes';

export class SomeNiceBuilderBuilder implements cog.Builder<with-dashes.SomeStruct> {
    private readonly internal: with-dashes.SomeStruct;

    constructor() {
        this.internal = with-dashes.defaultSomeStruct();
    }

    build(): with-dashes.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
