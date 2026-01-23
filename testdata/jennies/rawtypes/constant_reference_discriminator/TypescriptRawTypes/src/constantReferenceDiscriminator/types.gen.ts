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

export interface RowsLayoutUsingValue {
	kind: "RowsLayout";
	rowsLayoutProperty: string;
}

export const defaultRowsLayoutUsingValue = (): RowsLayoutUsingValue => ({
	kind: RowsLayoutKindType,
	rowsLayoutProperty: "",
});

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

export interface RowsLayoutWithoutValue {
	kind: "RowsLayout";
	rowsLayoutProperty: string;
}

export const defaultRowsLayoutWithoutValue = (): RowsLayoutWithoutValue => ({
	kind: RowsLayoutKindType,
	rowsLayoutProperty: "",
});

export const GridLayoutKindType = "GridLayout";

export const RowsLayoutKindType = "RowsLayout";
