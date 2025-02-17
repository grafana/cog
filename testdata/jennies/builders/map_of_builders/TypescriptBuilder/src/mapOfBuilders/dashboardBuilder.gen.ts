import * as cog from '../cog';
import * as mapOfBuilders from '../mapOfBuilders';

export class DashboardBuilder implements cog.Builder<mapOfBuilders.Dashboard> {
    protected readonly internal: mapOfBuilders.Dashboard;

    constructor() {
        this.internal = mapOfBuilders.defaultDashboard();
    }

    /**
     * Builds the object.
     */
    build(): mapOfBuilders.Dashboard {
        return this.internal;
    }

    panels(panels: Record<string, cog.Builder<mapOfBuilders.Panel>>): this {
        const panelsResource = (function() {
            let results1 = {};
            for (const key1 in panels) {
                const val1 = panels[key1];
                results1[key1] = val1.build();
            }
            return results1;
        }());
        this.internal.panels = panelsResource;
        return this;
    }
}

