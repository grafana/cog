import * as cog from '../cog';
import * as discriminatorWithoutOption from '../discriminatorWithoutOption';

export class NoShowFieldOptionBuilder implements cog.Builder<discriminatorWithoutOption.NoShowFieldOption> {
    protected readonly internal: discriminatorWithoutOption.NoShowFieldOption;

    constructor() {
        this.internal = discriminatorWithoutOption.defaultNoShowFieldOption();
    }

    /**
     * Builds the object.
     */
    build(): discriminatorWithoutOption.NoShowFieldOption {
        return this.internal;
    }

    text(text: string): this {
        this.internal.text = text;
        return this;
    }
}
