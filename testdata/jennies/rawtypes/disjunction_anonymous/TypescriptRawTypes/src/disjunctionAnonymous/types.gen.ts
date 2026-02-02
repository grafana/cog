export interface MyStruct {
	scalars: string | boolean | number | number;
	sameKind: "a" | "b" | "c";
	refs: StructA | StructB;
	mixed: StructA | string | number;
}

export const defaultMyStruct = (): MyStruct => ({
	scalars: "",
	sameKind: "a",
	refs: defaultStructA(),
	mixed: defaultStructA(),
});

export interface StructA {
	field: string;
}

export const defaultStructA = (): StructA => ({
	field: "",
});

export interface StructB {
	type: number;
}

export const defaultStructB = (): StructB => ({
	type: 0,
});

