export type LayoutWithValue = GridLayoutUsingValue | RowsLayoutUsingValue;

export const defaultLayoutWithValue = (): LayoutWithValue => (defaultGridLayoutUsingValue());

export interface GridLayoutUsingValue {
	kind: "GridLayout";
	gridLayoutProperty: string;
}

export const defaultGridLayoutUsingValue = (): GridLayoutUsingValue => ({
	kind: GridLayoutKindType,
	gridLayoutProperty: "",
});

// equalsGridLayoutUsingValue tests the equality of two `GridLayoutUsingValue` objects.
export const equalsGridLayoutUsingValue = (a: GridLayoutUsingValue, b: GridLayoutUsingValue): boolean => {
	if (a.kind !== b.kind) return false;
	if (a.gridLayoutProperty !== b.gridLayoutProperty) return false;
	return true;
};

export interface RowsLayoutUsingValue {
	kind: "RowsLayout";
	rowsLayoutProperty: string;
}

export const defaultRowsLayoutUsingValue = (): RowsLayoutUsingValue => ({
	kind: RowsLayoutKindType,
	rowsLayoutProperty: "",
});

// equalsRowsLayoutUsingValue tests the equality of two `RowsLayoutUsingValue` objects.
export const equalsRowsLayoutUsingValue = (a: RowsLayoutUsingValue, b: RowsLayoutUsingValue): boolean => {
	if (a.kind !== b.kind) return false;
	if (a.rowsLayoutProperty !== b.rowsLayoutProperty) return false;
	return true;
};

export type LayoutWithoutValue = GridLayoutWithoutValue | RowsLayoutWithoutValue;

export const defaultLayoutWithoutValue = (): LayoutWithoutValue => (defaultGridLayoutWithoutValue());

export interface GridLayoutWithoutValue {
	kind: "GridLayout";
	gridLayoutProperty: string;
}

export const defaultGridLayoutWithoutValue = (): GridLayoutWithoutValue => ({
	kind: GridLayoutKindType,
	gridLayoutProperty: "",
});

// equalsGridLayoutWithoutValue tests the equality of two `GridLayoutWithoutValue` objects.
export const equalsGridLayoutWithoutValue = (a: GridLayoutWithoutValue, b: GridLayoutWithoutValue): boolean => {
	if (a.kind !== b.kind) return false;
	if (a.gridLayoutProperty !== b.gridLayoutProperty) return false;
	return true;
};

export interface RowsLayoutWithoutValue {
	kind: "RowsLayout";
	rowsLayoutProperty: string;
}

export const defaultRowsLayoutWithoutValue = (): RowsLayoutWithoutValue => ({
	kind: RowsLayoutKindType,
	rowsLayoutProperty: "",
});

// equalsRowsLayoutWithoutValue tests the equality of two `RowsLayoutWithoutValue` objects.
export const equalsRowsLayoutWithoutValue = (a: RowsLayoutWithoutValue, b: RowsLayoutWithoutValue): boolean => {
	if (a.kind !== b.kind) return false;
	if (a.rowsLayoutProperty !== b.rowsLayoutProperty) return false;
	return true;
};

export const GridLayoutKindType = "GridLayout";

export const RowsLayoutKindType = "RowsLayout";

