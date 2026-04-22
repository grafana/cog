export enum Enum {
	ValueA = "ValueA",
	ValueB = "ValueB",
	ValueC = "ValueC",
}

export const defaultEnum = (): Enum => (Enum.ValueA);

export interface ParentStruct {
	myEnum: Enum;
}

export const defaultParentStruct = (): ParentStruct => ({
	myEnum: Enum.ValueA,
});

// equalsParentStruct tests the equality of two `ParentStruct` objects.
export const equalsParentStruct = (a: ParentStruct, b: ParentStruct): boolean => {
	if (a.myEnum !== b.myEnum) return false;
	return true;
};

export interface Struct {
	myValue: string;
	myEnum: Enum;
}

export const defaultStruct = (): Struct => ({
	myValue: "",
	myEnum: Enum.ValueA,
});

// equalsStruct tests the equality of two `Struct` objects.
export const equalsStruct = (a: Struct, b: Struct): boolean => {
	if (a.myValue !== b.myValue) return false;
	if (a.myEnum !== b.myEnum) return false;
	return true;
};

export interface StructA {
	myEnum: Enum.ValueA;
	other?: Enum.ValueA;
}

export const defaultStructA = (): StructA => ({
	myEnum: Enum.ValueA,
	other: Enum.ValueA,
});

// equalsStructA tests the equality of two `StructA` objects.
export const equalsStructA = (a: StructA, b: StructA): boolean => {
	if (a.myEnum !== b.myEnum) return false;
	if ((a.other === undefined) !== (b.other === undefined)) return false;
	if (a.other !== undefined) {
		if (a.other !== b.other!) return false;
	}
	return true;
};

export interface StructB {
	myEnum: Enum.ValueB;
	myValue: string;
}

export const defaultStructB = (): StructB => ({
	myEnum: Enum.ValueB,
	myValue: "",
});

// equalsStructB tests the equality of two `StructB` objects.
export const equalsStructB = (a: StructB, b: StructB): boolean => {
	if (a.myEnum !== b.myEnum) return false;
	if (a.myValue !== b.myValue) return false;
	return true;
};

