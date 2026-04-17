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

export interface SomeStructWithDefaultEnum {
	data: Record<StringEnumWithDefault, string>;
}

