import * as otherpkg from '../otherpkg';


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

export type RefToSomeStruct = SomeStruct;

export const defaultRefToSomeStruct = (): RefToSomeStruct => (defaultSomeStruct());

export type RefToSomeStructFromOtherPackage = otherpkg.SomeDistantStruct;

export const defaultRefToSomeStructFromOtherPackage = (): RefToSomeStructFromOtherPackage => (otherpkg.defaultSomeDistantStruct());

