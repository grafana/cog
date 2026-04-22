export interface Query {
	expr: string;
	instant?: boolean;
	_implementsDataqueryVariant(): void;
}

export const defaultQuery = (): Query => ({
	expr: "",
	_implementsDataqueryVariant: () => {},
});

// equalsQuery tests the equality of two `Query` objects.
export const equalsQuery = (a: Query, b: Query): boolean => {
	if (a.expr !== b.expr) return false;
	if ((a.instant === undefined) !== (b.instant === undefined)) return false;
	if (a.instant !== undefined) {
		if (a.instant !== b.instant!) return false;
	}
	return true;
};

