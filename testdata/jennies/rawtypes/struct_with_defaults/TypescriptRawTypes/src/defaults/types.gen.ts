export interface SomeStruct {
	fieldBool: boolean;
	fieldString: string;
	FieldStringWithConstantValue: "auto";
	FieldFloat32: number;
	FieldInt32: number;
}

export const defaultSomeStruct = (): SomeStruct => ({
	fieldBool: true,
	fieldString: "foo",
	FieldStringWithConstantValue: "auto",
	FieldFloat32: 42.42,
	FieldInt32: 42,
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (a.fieldBool !== b.fieldBool) return false;
	if (a.fieldString !== b.fieldString) return false;
	if (a.fieldStringWithConstantValue !== b.fieldStringWithConstantValue) return false;
	if (a.fieldFloat32 !== b.fieldFloat32) return false;
	if (a.fieldInt32 !== b.fieldInt32) return false;
	return true;
};

