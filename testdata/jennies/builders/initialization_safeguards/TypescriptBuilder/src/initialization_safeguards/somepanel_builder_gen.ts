import * as cog from '../cog';
import * as initialization_safeguards from '../initialization_safeguards';

export class SomePanelBuilder implements cog.Builder<initialization_safeguards.SomePanel> {
    private readonly internal: initialization_safeguards.SomePanel;

    constructor() {
        this.internal = initialization_safeguards.defaultSomePanel();
    }

    build(): initialization_safeguards.SomePanel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }

    showLegend(show: boolean): this {
        if (!this.internal.options) {
            this.internal.options = initialization_safeguards.defaultOptions();
        }
        if (!this.internal.options.legend) {
            this.internal.options.legend = {
	show: true,
};
        }
        this.internal.options.legend.show = show;
        return this;
    }
}
