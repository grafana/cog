import * as cog from '../cog';
import * as basicStructDefaults from '../basicStructDefaults';

export class SomeStructBuilder implements cog.Builder<basicStructDefaults.SomeStruct> {
    protected readonly internal: basicStructDefaults.SomeStruct;

    constructor() {
        this.internal = basicStructDefaults.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): basicStructDefaults.SomeStruct {
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

