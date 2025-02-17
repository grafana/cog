import * as cog from '../cog';
import * as initializationSafeguards from '../initializationSafeguards';

export class SomePanelBuilder implements cog.Builder<initializationSafeguards.SomePanel> {
    protected readonly internal: initializationSafeguards.SomePanel;

    constructor() {
        this.internal = initializationSafeguards.defaultSomePanel();
    }

    /**
     * Builds the object.
     */
    build(): initializationSafeguards.SomePanel {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }

    showLegend(show: boolean): this {
        if (!this.internal.options) {
            this.internal.options = initializationSafeguards.defaultOptions();
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

