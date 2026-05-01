export type DisjunctionWithoutDiscriminator = TypeA | TypeB;

export const defaultDisjunctionWithoutDiscriminator = (): DisjunctionWithoutDiscriminator => (defaultTypeA());

export interface TypeA {
	fieldA: string;
}

export const defaultTypeA = (): TypeA => ({
	fieldA: "",
});

// equalsTypeA tests the equality of two `TypeA` objects.
export const equalsTypeA = (a: TypeA, b: TypeA): boolean => {
	if (a.fieldA !== b.fieldA) return false;
	return true;
};

export interface TypeB {
	fieldB: number;
}

export const defaultTypeB = (): TypeB => ({
	fieldB: 0,
});

// equalsTypeB tests the equality of two `TypeB` objects.
export const equalsTypeB = (a: TypeB, b: TypeB): boolean => {
	if (a.fieldB !== b.fieldB) return false;
	return true;
};

