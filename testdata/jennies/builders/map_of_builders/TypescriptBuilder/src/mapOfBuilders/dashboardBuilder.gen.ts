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
            let results1: Record<string, mapOfBuilders.Panel> = {};
            for (const key1 in panels) {
                const val1 = panels[key1];
                results1[key1] = val1.build();
            }
            return results1;
        }());
        this.internal.panels = panelsResource;
        return this;
    }

    facets(facets: Record<string, cog.Builder<mapOfBuilders.FacetValue>[]>): this {
        const facetsResource = (function() {
            let results1: Record<string, mapOfBuilders.FacetValue[]> = {};
            for (const key1 in facets) {
                const val1 = facets[key1];
                results1[key1] = val1.map(builder2 => builder2.build());
            }
            return results1;
        }());
        this.internal.facets = facetsResource;
        return this;
    }
}

