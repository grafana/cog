import * as cog from '../cog';
import * as sandbox from '../sandbox';

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

    annotations(key: string,value: string): this {
        if (!this.internal.annotations) {
            this.internal.annotations = {};
        }
        this.internal.annotations[key] = value;
        return this;
    }
}

