import * as cog from '../cog';
import * as mapOfDisjunctions from '../mapOfDisjunctions';

export class LibraryPanelBuilder implements cog.Builder<mapOfDisjunctions.LibraryPanel> {
    protected readonly internal: mapOfDisjunctions.LibraryPanel;

    constructor() {
        this.internal = mapOfDisjunctions.defaultLibraryPanel();
        this.internal.kind = "Library";
    }

    /**
     * Builds the object.
     */
    build(): mapOfDisjunctions.LibraryPanel {
        return this.internal;
    }

    text(text: string): this {
        this.internal.text = text;
        return this;
    }
}
