import * as cog from '../cog';
import * as sandbox from '../sandbox';

export class SomeStructWithDefaultEnumBuilder implements cog.Builder<sandbox.SomeStructWithDefaultEnum> {
    protected readonly internal: sandbox.SomeStructWithDefaultEnum;

    constructor() {
        this.internal = sandbox.defaultSomeStructWithDefaultEnum();
    }

    /**
     * Builds the object.
     */
    build(): sandbox.SomeStructWithDefaultEnum {
        return this.internal;
    }

    data(key: sandbox.StringEnumWithDefault,value: string): this {
        if (!this.internal.data) {
            this.internal.data = {};
        }
        this.internal.data[key] = value;
        return this;
    }
}

