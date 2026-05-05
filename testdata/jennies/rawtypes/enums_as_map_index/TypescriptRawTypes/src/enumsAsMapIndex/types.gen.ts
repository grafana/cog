export enum StringEnum {
	A = "a",
	B = "b",
	C = "c",
}

export const defaultStringEnum = (): StringEnum => (StringEnum.A);

export enum StringEnumWithDefault {
	A = "a",
	B = "b",
	C = "c",
}

export const defaultStringEnumWithDefault = (): StringEnumWithDefault => (StringEnumWithDefault.A);

export interface SomeStruct {
	data: Record<StringEnum, string>;
}

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (Object.keys(a.data).length !== Object.keys(b.data).length) return false;
	for (const key1 in a.data) {
		if (a.data[key1] !== b.data[key1]) return false;
	}
	return true;
};

export interface SomeStructWithDefaultEnum {
	data: Record<StringEnumWithDefault, string>;
}

// equalsSomeStructWithDefaultEnum tests the equality of two `SomeStructWithDefaultEnum` objects.
export const equalsSomeStructWithDefaultEnum = (a: SomeStructWithDefaultEnum, b: SomeStructWithDefaultEnum): boolean => {
	if (Object.keys(a.data).length !== Object.keys(b.data).length) return false;
	for (const key1 in a.data) {
		if (a.data[key1] !== b.data[key1]) return false;
	}
	return true;
};

