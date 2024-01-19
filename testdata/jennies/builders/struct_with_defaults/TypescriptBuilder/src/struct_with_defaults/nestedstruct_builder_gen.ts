import * as cog from '../cog';
import * as struct_with_defaults from '../struct_with_defaults';

export class NestedStructBuilder implements cog.Builder<struct_with_defaults.NestedStruct> {
    private readonly internal: struct_with_defaults.NestedStruct;

    constructor() {
        this.internal = struct_with_defaults.defaultNestedStruct();
    }

    build(): struct_with_defaults.NestedStruct {
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
