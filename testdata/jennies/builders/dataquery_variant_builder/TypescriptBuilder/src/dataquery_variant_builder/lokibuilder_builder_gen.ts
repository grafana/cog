import * as cog from '../cog';
import * as dataquery_variant_builder from '../dataquery_variant_builder';

export class LokiBuilderBuilder implements cog.Builder<cog.Dataquery> {
    private readonly internal: dataquery_variant_builder.Loki;

    constructor() {
        this.internal = dataquery_variant_builder.defaultLoki();
    }

    build(): dataquery_variant_builder.Loki {
        return this.internal;
    }

    expr(expr: string): this {
        this.internal.expr = expr;
        return this;
    }
}
