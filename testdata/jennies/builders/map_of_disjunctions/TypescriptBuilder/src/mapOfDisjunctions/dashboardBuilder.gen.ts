import * as cog from '../cog';
import * as mapOfDisjunctions from '../mapOfDisjunctions';

export class DashboardBuilder implements cog.Builder<mapOfDisjunctions.Dashboard> {
    protected readonly internal: mapOfDisjunctions.Dashboard;

    constructor() {
        this.internal = mapOfDisjunctions.defaultDashboard();
    }

    /**
     * Builds the object.
     */
    build(): mapOfDisjunctions.Dashboard {
        return this.internal;
    }

    panels(panels: Record<string, cog.Builder<mapOfDisjunctions.Element>>): this {
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
