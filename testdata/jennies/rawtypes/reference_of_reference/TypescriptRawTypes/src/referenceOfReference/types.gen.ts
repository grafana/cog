export interface MyStruct {
	field?: OtherStruct;
}

export const defaultMyStruct = (): MyStruct => ({
});

// equalsMyStruct tests the equality of two `MyStruct` objects.
export const equalsMyStruct = (a: MyStruct, b: MyStruct): boolean => {
	if ((a.field === undefined) !== (b.field === undefined)) return false;
	if (a.field !== undefined) {
		if (!equalsOtherStruct(a.field, b.field!)) return false;
	}
	return true;
};

export type OtherStruct = AnotherStruct;

export const defaultOtherStruct = (): OtherStruct => (defaultAnotherStruct());

export interface AnotherStruct {
	a: string;
}

export const defaultAnotherStruct = (): AnotherStruct => ({
	a: "",
});

// equalsAnotherStruct tests the equality of two `AnotherStruct` objects.
export const equalsAnotherStruct = (a: AnotherStruct, b: AnotherStruct): boolean => {
	if (a.a !== b.a) return false;
	return true;
};

