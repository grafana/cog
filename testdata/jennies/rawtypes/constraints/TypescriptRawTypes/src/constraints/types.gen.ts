export interface SomeStruct {
	id: number;
	maybeId?: number;
	title: string;
	refStruct?: refStruct;
}

export const defaultSomeStruct = (): SomeStruct => ({
	id: 0,
	title: "",
});

export interface refStruct {
	labels: Record<string, string>;
	tags: string[];
}

export const defaultRefStruct = (): refStruct => ({
	labels: {},
	tags: [],
});

