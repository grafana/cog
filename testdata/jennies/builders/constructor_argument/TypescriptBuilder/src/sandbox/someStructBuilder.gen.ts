import * as cog from '../cog';
import * as sandbox from '../sandbox';

export class SomeStructBuilder implements cog.Builder<sandbox.SomeStruct> {
    protected readonly internal: sandbox.SomeStruct;

    constructor(title: string) {
        this.internal = sandbox.defaultSomeStruct();
        this.internal.title = title;
    }

    /**
     * Builds the object.
     */
    build(): sandbox.SomeStruct {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}

