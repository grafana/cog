import * as cog from '../cog';
import * as mapOfBuilders from '../mapOfBuilders';

export class PanelBuilder implements cog.Builder<mapOfBuilders.Panel> {
    protected readonly internal: mapOfBuilders.Panel;

    constructor() {
        this.internal = mapOfBuilders.defaultPanel();
    }

    /**
     * Builds the object.
     */
    build(): mapOfBuilders.Panel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}

