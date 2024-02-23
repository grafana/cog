import * as cog from '../cog';
import * as nullable_map_assignment from '../nullable_map_assignment';

export class SomeStructBuilder implements cog.Builder<nullable_map_assignment.SomeStruct> {
    private readonly internal: nullable_map_assignment.SomeStruct;

    constructor() {
        this.internal = nullable_map_assignment.defaultSomeStruct();
    }

    build(): nullable_map_assignment.SomeStruct {
        return this.internal;
    }

    config(config: Record<string, string>): this {
        this.internal.config = config;
        return this;
    }
}
