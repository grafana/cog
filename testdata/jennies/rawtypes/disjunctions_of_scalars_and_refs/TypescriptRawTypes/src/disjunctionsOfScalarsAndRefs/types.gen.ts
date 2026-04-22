export type DisjunctionOfScalarsAndRefs = "a" | boolean | string[] | MyRefA | MyRefB;

export const defaultDisjunctionOfScalarsAndRefs = (): DisjunctionOfScalarsAndRefs => ("a");

export interface MyRefA {
	foo: string;
}

export const defaultMyRefA = (): MyRefA => ({
	foo: "",
});

// equalsMyRefA tests the equality of two `MyRefA` objects.
export const equalsMyRefA = (a: MyRefA, b: MyRefA): boolean => {
	if (a.foo !== b.foo) return false;
	return true;
};

export interface MyRefB {
	bar: number;
}

export const defaultMyRefB = (): MyRefB => ({
	bar: 0,
});

// equalsMyRefB tests the equality of two `MyRefB` objects.
export const equalsMyRefB = (a: MyRefB, b: MyRefB): boolean => {
	if (a.bar !== b.bar) return false;
	return true;
};

