import * as cog from '../cog';
import * as sandbox from '../sandbox';

/**
 * @deprecated This builder is deprecated. Don't use. Please.
 */
export class SomeStructBuilder implements cog.Builder<sandbox.SomeStruct> {
    protected readonly internal: sandbox.SomeStruct;

    constructor() {
        this.internal = sandbox.defaultSomeStruct();
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

