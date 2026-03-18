export type objTime = string;

export const defaultObjTime = (): objTime => ("");

export interface objWithTimeField {
	registeredAt: string;
	duration: string;
}

export const defaultObjWithTimeField = (): objWithTimeField => ({
	registeredAt: "",
	duration: "",
});

