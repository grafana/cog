import * as cog from '../cog';
import * as mapOfDisjunctions from '../mapOfDisjunctions';

export class ElementBuilder implements cog.Builder<mapOfDisjunctions.Element> {
    protected readonly internal: mapOfDisjunctions.Element;

    constructor() {
        this.internal = mapOfDisjunctions.defaultElement();
    }

    /**
     * Builds the object.
     */
    build(): mapOfDisjunctions.Element {
        return this.internal;
    }

    panel(panel: cog.Builder<mapOfDisjunctions.Panel>): this {
        const panelResource = panel.build();
        this.internal.Panel = panelResource;
        return this;
    }

    libraryPanel(libraryPanel: cog.Builder<mapOfDisjunctions.LibraryPanel>): this {
        const libraryPanelResource = libraryPanel.build();
        this.internal.LibraryPanel = libraryPanelResource;
        return this;
    }
}
