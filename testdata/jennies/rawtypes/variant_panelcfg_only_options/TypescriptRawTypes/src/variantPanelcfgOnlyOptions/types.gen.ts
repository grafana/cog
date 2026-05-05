export interface Options {
	content: string;
}

export const defaultOptions = (): Options => ({
	content: "",
});

// equalsOptions tests the equality of two `Options` objects.
export const equalsOptions = (a: Options, b: Options): boolean => {
	if (a.content !== b.content) return false;
	return true;
};

