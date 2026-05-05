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

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (a.id !== b.id) return false;
	if ((a.maybeId === undefined) !== (b.maybeId === undefined)) return false;
	if (a.maybeId !== undefined) {
		if (a.maybeId !== b.maybeId!) return false;
	}
	if (a.greaterThanZero !== b.greaterThanZero) return false;
	if (a.negative !== b.negative) return false;
	if (a.title !== b.title) return false;
	if (Object.keys(a.labels).length !== Object.keys(b.labels).length) return false;
	for (const key1 in a.labels) {
		if (a.labels[key1] !== b.labels[key1]) return false;
	}
	if (a.tags.length !== b.tags.length) return false;
	for (let i1 = 0; i1 < a.tags.length; i1++) {
		if (a.tags[i1] !== b.tags[i1]) return false;
	}
	if (a.regex !== b.regex) return false;
	if (a.negativeRegex !== b.negativeRegex) return false;
	if (a.minMaxList.length !== b.minMaxList.length) return false;
	for (let i1 = 0; i1 < a.minMaxList.length; i1++) {
		if (a.minMaxList[i1] !== b.minMaxList[i1]) return false;
	}
	if (a.uniqueList.length !== b.uniqueList.length) return false;
	for (let i1 = 0; i1 < a.uniqueList.length; i1++) {
		if (a.uniqueList[i1] !== b.uniqueList[i1]) return false;
	}
	if (a.fullConstraintList.length !== b.fullConstraintList.length) return false;
	for (let i1 = 0; i1 < a.fullConstraintList.length; i1++) {
		if (a.fullConstraintList[i1] !== b.fullConstraintList[i1]) return false;
	}
	return true;
};

