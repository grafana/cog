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

export type MapOfStringToRef = Record<string, SomeStruct>;

export const defaultMapOfStringToRef = (): MapOfStringToRef => ({});

export type MapOfStringToMapOfStringToBool = Record<string, Record<string, boolean>>;

export const defaultMapOfStringToMapOfStringToBool = (): MapOfStringToMapOfStringToBool => ({});

