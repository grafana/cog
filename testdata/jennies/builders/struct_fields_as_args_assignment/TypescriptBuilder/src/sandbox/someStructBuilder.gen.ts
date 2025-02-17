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

    time(from: string,to: string): this {
        if (!this.internal.time) {
            this.internal.time = {
	from: "now-6h",
	to: "now",
};
        }
        this.internal.time.from = from;
        this.internal.time.to = to;
        return this;
    }
}

