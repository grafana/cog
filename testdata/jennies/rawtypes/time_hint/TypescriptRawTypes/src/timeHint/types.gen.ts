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

// equalsobjWithTimeField tests the equality of two `objWithTimeField` objects.
export const equalsobjWithTimeField = (a: objWithTimeField, b: objWithTimeField): boolean => {
	if (a.registeredAt !== b.registeredAt) return false;
	if (a.duration !== b.duration) return false;
	return true;
};

