import * as cog from '../cog';
import * as anonymousStruct from '../anonymousStruct';

export class SomeStructBuilder implements cog.Builder<anonymousStruct.SomeStruct> {
    protected readonly internal: anonymousStruct.SomeStruct;

    constructor() {
        this.internal = anonymousStruct.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): anonymousStruct.SomeStruct {
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

