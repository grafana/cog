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

    editable(): this {
        this.internal.editable = true;
        return this;
    }

    readonly(): this {
        this.internal.editable = false;
        return this;
    }

    autoRefresh(): this {
        this.internal.autoRefresh = true;
        return this;
    }

    noAutoRefresh(): this {
        this.internal.autoRefresh = false;
        return this;
    }
}

