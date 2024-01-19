import * as cog from '../cog';
import * as basic_struct_defaults from '../basic_struct_defaults';

export class SomeStructBuilder implements cog.Builder<basic_struct_defaults.SomeStruct> {
    private readonly internal: basic_struct_defaults.SomeStruct;

    constructor() {
        this.internal = basic_struct_defaults.defaultSomeStruct();
    }

    build(): basic_struct_defaults.SomeStruct {
        return this.internal;
    }

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

    liveNow(liveNow: boolean): this {
        this.internal.liveNow = liveNow;
        return this;
    }
}
