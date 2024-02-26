import * as otherpkg from '../otherpkg';


export interface SomeStruct {
	FieldAny: any;
}

export const defaultSomeStruct = (): SomeStruct => ({
	FieldAny: {},
});

export type RefToSomeStruct = SomeStruct;

export const defaultRefToSomeStruct = (): RefToSomeStruct => (defaultSomeStruct());

export type RefToSomeStructFromOtherPackage = otherpkg.SomeDistantStruct;

export const defaultRefToSomeStructFromOtherPackage = (): RefToSomeStructFromOtherPackage => (otherpkg.default());

