export interface SomeStruct {
	id: number;
	maybeId?: number;
	greaterThanZero: number;
	negative: number;
	title: string;
	labels: Record<string, string>;
	tags: string[];
	regex: string;
	negativeRegex: string;
	minMaxList: string[];
	uniqueList: string[];
	fullConstraintList: number[];
}

export const defaultSomeStruct = (): SomeStruct => ({
	id: 0,
	greaterThanZero: 0,
	negative: 0,
	title: "",
	labels: {},
	tags: [],
	regex: "",
	negativeRegex: "",
	minMaxList: [],
	uniqueList: [],
	fullConstraintList: [],
});

