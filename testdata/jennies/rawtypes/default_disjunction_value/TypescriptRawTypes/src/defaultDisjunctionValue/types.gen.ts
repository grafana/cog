export type DisjunctionClasses = ValueA | ValueB | ValueC;

export const defaultDisjunctionClasses = (): DisjunctionClasses => (defaultValueA());

export interface ValueA {
	type: "A";
	anArray: string[];
	otherRef: ValueB;
}

export const defaultValueA = (): ValueA => ({
	type: "A",
	anArray: [],
	otherRef: defaultValueB(),
});

export interface ValueB {
	type: "B";
	aMap: Record<string, number>;
	def: 1 | "a" | boolean;
}

export const defaultValueB = (): ValueB => ({
	type: "B",
	aMap: {},
	def: 1,
});

export interface ValueC {
	type: "C";
	other: number;
}

export const defaultValueC = (): ValueC => ({
	type: "C",
	other: 0,
});

export type DisjunctionConstants = "abc" | 1 | true;

export const defaultDisjunctionConstants = (): DisjunctionConstants => (1);
