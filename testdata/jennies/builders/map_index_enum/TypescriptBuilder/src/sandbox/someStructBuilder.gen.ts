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

    data(key: sandbox.StringEnum,value: string): this {
        if (!this.internal.data) {
            this.internal.data = {};
        }
        this.internal.data[key] = value;
        return this;
    }
}

