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

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (a.type !== b.type) return false;
	if (JSON.stringify(a.fieldAny) !== JSON.stringify(b.fieldAny)) return false;
	return true;
};

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

// equalsSomeOtherStruct tests the equality of two `SomeOtherStruct` objects.
export const equalsSomeOtherStruct = (a: SomeOtherStruct, b: SomeOtherStruct): boolean => {
	if (a.type !== b.type) return false;
	if (a.foo !== b.foo) return false;
	return true;
};

export interface YetAnotherStruct {
	Type: "yet-another-struct";
	Bar: number;
}

export const defaultYetAnotherStruct = (): YetAnotherStruct => ({
	Type: "yet-another-struct",
	Bar: 0,
});

// equalsYetAnotherStruct tests the equality of two `YetAnotherStruct` objects.
export const equalsYetAnotherStruct = (a: YetAnotherStruct, b: YetAnotherStruct): boolean => {
	if (a.type !== b.type) return false;
	if (a.bar !== b.bar) return false;
	return true;
};

export type SeveralRefs = SomeStruct | SomeOtherStruct | YetAnotherStruct;

export const defaultSeveralRefs = (): SeveralRefs => (defaultSomeStruct());

