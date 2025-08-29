export type DisjunctionOfScalarsAndRefs = "a" | boolean | string[] | MyRefA | MyRefB;

export const defaultDisjunctionOfScalarsAndRefs = (): DisjunctionOfScalarsAndRefs => ("a");

export interface MyRefA {
	foo: string;
}

export const defaultMyRefA = (): MyRefA => ({
	foo: "",
});

export interface MyRefB {
	bar: number;
}

export const defaultMyRefB = (): MyRefB => ({
	bar: 0,
});
