import * as cog from '../cog';
import * as anonymous_struct from '../anonymous_struct';

export class SomeStructBuilder implements cog.Builder<anonymous_struct.SomeStruct> {
    private readonly internal: anonymous_struct.SomeStruct;

    constructor() {
        this.internal = anonymous_struct.defaultSomeStruct();
    }

    build(): anonymous_struct.SomeStruct {
        return this.internal;
    }

    time(time: {
	from: string;
	to: string;
}): this {
        this.internal.time = time;
        return this;
    }
}
