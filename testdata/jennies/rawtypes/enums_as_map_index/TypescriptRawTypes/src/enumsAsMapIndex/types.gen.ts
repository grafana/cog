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

export const defaultSomeStruct = (): SomeStruct => ({
	data: {
        [StringEnum.A]: "",
        [StringEnum.B]: "",
        [StringEnum.C]: ""
    },
});

export interface SomeStructWithDefaultEnum {
	data: Record<StringEnumWithDefault, string>;
}

export const defaultSomeStructWithDefaultEnum = (): SomeStructWithDefaultEnum => ({
	data: {},
});

