export interface SomeStruct {
	id: number;
	maybeId?: number;
	greaterThanZero: number;
	title: string;
	labels: Record<string, string>;
	tags: string[];
}

export const defaultSomeStruct = (): SomeStruct => ({
	id: 0,
	greaterThanZero: 0,
	title: "",
	labels: {},
	tags: [],
});

