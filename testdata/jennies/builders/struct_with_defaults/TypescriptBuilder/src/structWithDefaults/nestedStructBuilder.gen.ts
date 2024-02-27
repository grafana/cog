import * as cog from '../cog';
import * as structWithDefaults from '../structWithDefaults';

export class NestedStructBuilder implements cog.Builder<structWithDefaults.NestedStruct> {
    private readonly internal: structWithDefaults.NestedStruct;

    constructor() {
        this.internal = structWithDefaults.defaultNestedStruct();
    }

    build(): structWithDefaults.NestedStruct {
        return this.internal;
    }

    stringVal(stringVal: string): this {
        this.internal.stringVal = stringVal;
        return this;
    }

    intVal(intVal: number): this {
        this.internal.intVal = intVal;
        return this;
    }
}
