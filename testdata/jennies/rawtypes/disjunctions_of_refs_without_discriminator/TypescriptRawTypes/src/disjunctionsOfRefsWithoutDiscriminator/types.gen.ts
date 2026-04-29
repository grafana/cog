export type DisjunctionWithoutDiscriminator = TypeA | TypeB;

export const defaultDisjunctionWithoutDiscriminator = (): DisjunctionWithoutDiscriminator => (defaultTypeA());

export interface TypeA {
	fieldA: string;
}

export const defaultTypeA = (): TypeA => ({
	fieldA: "",
});

export interface TypeB {
	fieldB: number;
}

export const defaultTypeB = (): TypeB => ({
	fieldB: 0,
});

