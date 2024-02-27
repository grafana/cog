export interface someStruct {
	FieldAny: any;
}

export const defaultSomeStruct = (): someStruct => ({
	FieldAny: {},
});

// Refresh rate or disabled.
export type RefreshRate = string | boolean;

export const defaultRefreshRate = (): RefreshRate => ("");

