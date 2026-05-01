// List of tags, maybe?
export type ArrayOfStrings = string[];

export const defaultArrayOfStrings = (): ArrayOfStrings => ([]);

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

export type ArrayOfRefs = someStruct[];

export const defaultArrayOfRefs = (): ArrayOfRefs => ([]);

export type ArrayOfArrayOfNumbers = number[][];

export const defaultArrayOfArrayOfNumbers = (): ArrayOfArrayOfNumbers => ([]);

