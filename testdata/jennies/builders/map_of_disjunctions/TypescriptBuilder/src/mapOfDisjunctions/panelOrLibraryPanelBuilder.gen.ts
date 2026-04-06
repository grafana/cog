import * as cog from '../cog';
import * as mapOfDisjunctions from '../mapOfDisjunctions';

export class PanelOrLibraryPanelBuilder implements cog.Builder<mapOfDisjunctions.PanelOrLibraryPanel> {
    protected readonly internal: mapOfDisjunctions.PanelOrLibraryPanel;

    constructor() {
        this.internal = mapOfDisjunctions.defaultPanelOrLibraryPanel();
    }

    /**
     * Builds the object.
     */
    build(): mapOfDisjunctions.PanelOrLibraryPanel {
        return this.internal;
    }

    panel(panel: cog.Builder<mapOfDisjunctions.Panel>): this {
        const panelResource = panel.build();
        this.internal.panel = panelResource;
        return this;
    }

    libraryPanel(libraryPanel: cog.Builder<mapOfDisjunctions.LibraryPanel>): this {
        const libraryPanelResource = libraryPanel.build();
        this.internal.libraryPanel = libraryPanelResource;
        return this;
    }
}

