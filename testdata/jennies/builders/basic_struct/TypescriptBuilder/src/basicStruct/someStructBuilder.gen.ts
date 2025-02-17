import * as cog from '../cog';
import * as basicStruct from '../basicStruct';

// SomeStruct, to hold data.
export class SomeStructBuilder implements cog.Builder<basicStruct.SomeStruct> {
    protected readonly internal: basicStruct.SomeStruct;

    constructor() {
        this.internal = basicStruct.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): basicStruct.SomeStruct {
        return this.internal;
    }

    // id identifies something. Weird, right?
    id(id: number): this {
        this.internal.id = id;
        return this;
    }

    uid(uid: string): this {
        this.internal.uid = uid;
        return this;
    }

    tags(tags: string[]): this {
        this.internal.tags = tags;
        return this;
    }

    // This thing could be live.
    // Or maybe not.
    liveNow(liveNow: boolean): this {
        this.internal.liveNow = liveNow;
        return this;
    }
}

