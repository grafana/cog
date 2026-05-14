/**
* @deprecated: This object is deprecated, use NewStruct instead.
*/
export interface SomeStruct {
	FieldString: string;
}

export const defaultSomeStruct = (): SomeStruct => ({
	FieldString: "",
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (a.fieldString !== b.fieldString) return false;
	return true;
};

