import * as cog from '../cog';
import * as constructorInitializations from '../constructorInitializations';

export class SomePanelBuilder implements cog.Builder<constructorInitializations.SomePanel> {
    private readonly internal: constructorInitializations.SomePanel;

    constructor() {
        this.internal = constructorInitializations.defaultSomePanel();
        this.internal.type = "panel_type";
        this.internal.cursor = constructorInitializations.CursorMode.Tooltip;
    }

    build(): constructorInitializations.SomePanel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }
}
