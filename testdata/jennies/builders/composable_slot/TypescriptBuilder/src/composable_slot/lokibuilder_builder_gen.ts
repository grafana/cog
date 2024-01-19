import * as cog from '../cog';
import * as composable_slot from '../composable_slot';

export class LokiBuilderBuilder implements cog.Builder<composable_slot.Dashboard> {
    private readonly internal: composable_slot.Dashboard;

    constructor() {
        this.internal = composable_slot.defaultDashboard();
    }

    build(): composable_slot.Dashboard {
        return this.internal;
    }

    target(target: cog.Builder<cog.Dataquery>): this {
        const targetResource = target.build();
        this.internal.target = targetResource;
        return this;
    }

    targets(targets: cog.Builder<cog.Dataquery>[]): this {
        const targetsResources = targets.map(builder => builder.build());
        this.internal.targets = targetsResources;
        return this;
    }
}
