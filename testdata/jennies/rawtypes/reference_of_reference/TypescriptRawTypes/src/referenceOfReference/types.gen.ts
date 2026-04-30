export interface MyStruct {
	field?: OtherStruct;
}

export const defaultMyStruct = (): MyStruct => ({
});

export type OtherStruct = AnotherStruct;

export const defaultOtherStruct = (): OtherStruct => (defaultAnotherStruct());

export interface AnotherStruct {
	a: string;
}

export const defaultAnotherStruct = (): AnotherStruct => ({
	a: "",
});

