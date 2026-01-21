import * as cog from '../cog';
import * as discriminatorWithoutOption from '../discriminatorWithoutOption';

export class ShowFieldOptionBuilder implements cog.Builder<discriminatorWithoutOption.ShowFieldOption> {
    protected readonly internal: discriminatorWithoutOption.ShowFieldOption;

    constructor() {
        this.internal = discriminatorWithoutOption.defaultShowFieldOption();
    }

    /**
     * Builds the object.
     */
    build(): discriminatorWithoutOption.ShowFieldOption {
        return this.internal;
    }

    field(field: discriminatorWithoutOption.AnEnum): this {
        this.internal.field = field;
        return this;
    }

    text(text: string): this {
        this.internal.text = text;
        return this;
    }
}
