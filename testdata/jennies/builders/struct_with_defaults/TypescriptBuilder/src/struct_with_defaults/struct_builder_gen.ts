import * as cog from '../cog';
import * as struct_with_defaults from '../struct_with_defaults';

export class StructBuilder implements cog.Builder<struct_with_defaults.Struct> {
    private readonly internal: struct_with_defaults.Struct;

    constructor() {
        this.internal = struct_with_defaults.defaultStruct();
    }

    build(): struct_with_defaults.Struct {
        return this.internal;
    }

    allFields(allFields: cog.Builder<struct_with_defaults.NestedStruct>): this {
        const allFieldsResource = allFields.build();
        this.internal.allFields = allFieldsResource;
        return this;
    }

    partialFields(partialFields: cog.Builder<struct_with_defaults.NestedStruct>): this {
        const partialFieldsResource = partialFields.build();
        this.internal.partialFields = partialFieldsResource;
        return this;
    }

    emptyFields(emptyFields: cog.Builder<struct_with_defaults.NestedStruct>): this {
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
