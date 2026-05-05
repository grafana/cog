// String to... something.
export type MapOfStringToAny = Record<string, any>;

export const defaultMapOfStringToAny = (): MapOfStringToAny => ({});

export type MapOfStringToString = Record<string, string>;

export const defaultMapOfStringToString = (): MapOfStringToString => ({});

export interface SomeStruct {
	FieldAny: any;
}

export const defaultSomeStruct = (): SomeStruct => ({
	FieldAny: {},
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (JSON.stringify(a.fieldAny) !== JSON.stringify(b.fieldAny)) return false;
	return true;
};

export type MapOfStringToRef = Record<string, SomeStruct>;

export const defaultMapOfStringToRef = (): MapOfStringToRef => ({});

export type MapOfStringToMapOfStringToBool = Record<string, Record<string, boolean>>;

export const defaultMapOfStringToMapOfStringToBool = (): MapOfStringToMapOfStringToBool => ({});

