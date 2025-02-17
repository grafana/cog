import * as cog from '../cog';
import * as structWithDefaults from '../structWithDefaults';

export class StructBuilder implements cog.Builder<structWithDefaults.Struct> {
    protected readonly internal: structWithDefaults.Struct;

    constructor() {
        this.internal = structWithDefaults.defaultStruct();
    }

    /**
     * Builds the object.
     */
    build(): structWithDefaults.Struct {
        return this.internal;
    }

    allFields(allFields: cog.Builder<structWithDefaults.NestedStruct>): this {
        const allFieldsResource = allFields.build();
        this.internal.allFields = allFieldsResource;
        return this;
    }

    partialFields(partialFields: cog.Builder<structWithDefaults.NestedStruct>): this {
        const partialFieldsResource = partialFields.build();
        this.internal.partialFields = partialFieldsResource;
        return this;
    }

    emptyFields(emptyFields: cog.Builder<structWithDefaults.NestedStruct>): this {
        const emptyFieldsResource = emptyFields.build();
        this.internal.emptyFields = emptyFieldsResource;
        return this;
    }

    complexField(complexField: {
	uid: string;
	nested: {
		nestedVal: string;
	};
	array: string[];
}): this {
        this.internal.complexField = complexField;
        return this;
    }

    partialComplexField(partialComplexField: {
	uid: string;
	intVal: number;
}): this {
        this.internal.partialComplexField = partialComplexField;
        return this;
    }
}

