import * as cog from '../cog';
import * as properties from '../properties';

export class SomeStructBuilder implements cog.Builder<properties.SomeStruct> {
    protected readonly internal: properties.SomeStruct;
    private someBuilderProperty: string = "";

    constructor() {
        this.internal = properties.defaultSomeStruct();
    }

    /**
     * Builds the object.
     */
    build(): properties.SomeStruct {
        return this.internal;
    }

    id(id: number): this {
        this.internal.id = id;
        return this;
    }
}

