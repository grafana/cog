export interface someStruct {
	FieldAny: any;
}

export const defaultSomeStruct = (): someStruct => ({
	FieldAny: {},
});

// equalssomeStruct tests the equality of two `someStruct` objects.
export const equalssomeStruct = (a: someStruct, b: someStruct): boolean => {
	if (JSON.stringify(a.fieldAny) !== JSON.stringify(b.fieldAny)) return false;
	return true;
};

// Refresh rate or disabled.
export type RefreshRate = string | boolean;

export const defaultRefreshRate = (): RefreshRate => ("");

