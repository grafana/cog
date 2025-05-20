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

export interface Struct {
	myValue: string;
	myEnum: Enum;
}

export const defaultStruct = (): Struct => ({
	myValue: "",
	myEnum: Enum.ValueA,
});

export interface StructA {
	myEnum: Enum.ValueA;
}

export const defaultStructA = (): StructA => ({
	myEnum: Enum.ValueA,
});

export interface StructB {
	myEnum: Enum.ValueB;
	myValue: string;
}

export const defaultStructB = (): StructB => ({
	myEnum: Enum.ValueB,
	myValue: "",
});

