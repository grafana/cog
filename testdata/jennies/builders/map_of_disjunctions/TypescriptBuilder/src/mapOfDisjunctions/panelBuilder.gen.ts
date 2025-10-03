import * as cog from '../cog';
import * as mapOfDisjunctions from '../mapOfDisjunctions';

export class PanelBuilder implements cog.Builder<mapOfDisjunctions.Panel> {
    protected readonly internal: mapOfDisjunctions.Panel;

    constructor() {
        this.internal = mapOfDisjunctions.defaultPanel();
        this.internal.kind = "Panel";
    }

    /**
     * Builds the object.
     */
    build(): mapOfDisjunctions.Panel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
