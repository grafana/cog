import * as cog from '../cog';
import * as dataqueryVariantBuilder from '../dataqueryVariantBuilder';

export class LokiBuilderBuilder implements cog.Builder<cog.Dataquery> {
    protected readonly internal: dataqueryVariantBuilder.Loki;

    constructor() {
        this.internal = dataqueryVariantBuilder.defaultLoki();
    }

    /**
     * Builds the object.
     */
    build(): dataqueryVariantBuilder.Loki {
        return this.internal;
    }

    expr(expr: string): this {
        this.internal.expr = expr;
        return this;
    }
}

