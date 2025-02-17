import * as cog from '../cog';
import * as composableSlot from '../composableSlot';

export class LokiBuilderBuilder implements cog.Builder<composableSlot.Dashboard> {
    protected readonly internal: composableSlot.Dashboard;

    constructor() {
        this.internal = composableSlot.defaultDashboard();
    }

    /**
     * Builds the object.
     */
    build(): composableSlot.Dashboard {
        return this.internal;
    }

    target(target: cog.Builder<cog.Dataquery>): this {
        const targetResource = target.build();
        this.internal.target = targetResource;
        return this;
    }

    targets(targets: cog.Builder<cog.Dataquery>[]): this {
        const targetsResources = targets.map(builder1 => builder1.build());
        this.internal.targets = targetsResources;
        return this;
    }
}

