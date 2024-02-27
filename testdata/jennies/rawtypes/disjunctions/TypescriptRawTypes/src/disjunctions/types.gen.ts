// Refresh rate or disabled.
export type RefreshRate = string | boolean;

export const defaultRefreshRate = (): RefreshRate => ("");

export type StringOrNull = string | null;

export const defaultStringOrNull = (): StringOrNull => ("");

export interface SomeStruct {
	Type: "some-struct";
	FieldAny: any;
}

export const defaultSomeStruct = (): SomeStruct => ({
	Type: "some-struct",
	FieldAny: {},
});

export type BoolOrRef = boolean | SomeStruct;

export const defaultBoolOrRef = (): BoolOrRef => (false);

export interface SomeOtherStruct {
	Type: "some-other-struct";
	Foo: string;
}

export const defaultSomeOtherStruct = (): SomeOtherStruct => ({
	Type: "some-other-struct",
	Foo: "",
});

export interface YetAnotherStruct {
	Type: "yet-another-struct";
	Bar: number;
}

export const defaultYetAnotherStruct = (): YetAnotherStruct => ({
	Type: "yet-another-struct",
	Bar: 0,
});

export type SeveralRefs = SomeStruct | SomeOtherStruct | YetAnotherStruct;

export const defaultSeveralRefs = (): SeveralRefs => (defaultSomeStruct());

