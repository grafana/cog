import * as cog from '../cog';
import * as mapOfBuilders from '../mapOfBuilders';

export class FacetValueBuilder implements cog.Builder<mapOfBuilders.FacetValue> {
    protected readonly internal: mapOfBuilders.FacetValue;

    constructor() {
        this.internal = mapOfBuilders.defaultFacetValue();
    }

    /**
     * Builds the object.
     */
    build(): mapOfBuilders.FacetValue {
        return this.internal;
    }

    value(value: string): this {
        this.internal.value = value;
        return this;
    }

    count(count: number): this {
        this.internal.count = count;
        return this;
    }
}

