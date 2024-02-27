// List of tags, maybe?
export type ArrayOfStrings = string[];

export const defaultArrayOfStrings = (): ArrayOfStrings => ([]);

export interface someStruct {
	FieldAny: any;
}

export const defaultSomeStruct = (): someStruct => ({
	FieldAny: {},
});

export type ArrayOfRefs = someStruct[];

export const defaultArrayOfRefs = (): ArrayOfRefs => ([]);

export type ArrayOfArrayOfNumbers = number[][];

export const defaultArrayOfArrayOfNumbers = (): ArrayOfArrayOfNumbers => ([]);

