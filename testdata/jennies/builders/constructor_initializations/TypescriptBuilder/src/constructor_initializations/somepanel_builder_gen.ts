import * as cog from '../cog';
import * as constructor_initializations from '../constructor_initializations';

export class SomePanelBuilder implements cog.Builder<constructor_initializations.SomePanel> {
    private readonly internal: constructor_initializations.SomePanel;

    constructor() {
        this.internal = constructor_initializations.defaultSomePanel();
        this.internal.type = "panel_type";
        this.internal.cursor = constructor_initializations.CursorMode.Tooltip;
    }

    build(): constructor_initializations.SomePanel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
