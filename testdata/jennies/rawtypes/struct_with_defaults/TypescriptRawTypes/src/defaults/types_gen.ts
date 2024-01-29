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

