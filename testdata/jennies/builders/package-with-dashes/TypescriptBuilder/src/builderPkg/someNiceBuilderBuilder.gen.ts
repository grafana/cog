import * as cog from '../cog';
import * as withDashes from '../withDashes';

export class SomeNiceBuilderBuilder implements cog.Builder<withDashes.SomeStruct> {
    private readonly internal: withDashes.SomeStruct;

    constructor() {
        this.internal = withDashes.defaultSomeStruct();
    }

    build(): withDashes.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
