import * as cog from '../cog';
import * as basic_struct from '../basic_struct';

// SomeStruct, to hold data.
export class SomeStructBuilder implements cog.Builder<basic_struct.SomeStruct> {
    private readonly internal: basic_struct.SomeStruct;

    constructor() {
        this.internal = basic_struct.defaultSomeStruct();
    }

    build(): basic_struct.SomeStruct {
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
