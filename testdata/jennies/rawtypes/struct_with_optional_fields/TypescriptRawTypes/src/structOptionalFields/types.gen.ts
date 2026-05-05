export interface SomeStruct {
	FieldRef?: SomeOtherStruct;
	FieldString?: string;
	Operator?: ">" | "<";
	FieldArrayOfStrings?: string[];
	FieldAnonymousStruct?: {
		FieldAny: any;
	};
}

export const defaultSomeStruct = (): SomeStruct => ({
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if ((a.fieldRef === undefined) !== (b.fieldRef === undefined)) return false;
	if (a.fieldRef !== undefined) {
		if (!equalsSomeOtherStruct(a.fieldRef, b.fieldRef!)) return false;
	}
	if ((a.fieldString === undefined) !== (b.fieldString === undefined)) return false;
	if (a.fieldString !== undefined) {
		if (a.fieldString !== b.fieldString!) return false;
	}
	if ((a.operator === undefined) !== (b.operator === undefined)) return false;
	if (a.operator !== undefined) {
		if (a.operator !== b.operator!) return false;
	}
	if ((a.fieldArrayOfStrings === undefined) !== (b.fieldArrayOfStrings === undefined)) return false;
	if (a.fieldArrayOfStrings !== undefined) {
		if (a.fieldArrayOfStrings.length !== b.fieldArrayOfStrings!.length) return false;
		for (let i2 = 0; i2 < a.fieldArrayOfStrings.length; i2++) {
			if (a.fieldArrayOfStrings[i2] !== b.fieldArrayOfStrings![i2]) return false;
		}
	}
	if ((a.fieldAnonymousStruct === undefined) !== (b.fieldAnonymousStruct === undefined)) return false;
	if (a.fieldAnonymousStruct !== undefined) {
		if (!equalsUnknown(a.fieldAnonymousStruct, b.fieldAnonymousStruct!)) return false;
	}
	return true;
};

export interface SomeOtherStruct {
	FieldAny: any;
}

export const defaultSomeOtherStruct = (): SomeOtherStruct => ({
	FieldAny: {},
});

// equalsSomeOtherStruct tests the equality of two `SomeOtherStruct` objects.
export const equalsSomeOtherStruct = (a: SomeOtherStruct, b: SomeOtherStruct): boolean => {
	if (JSON.stringify(a.fieldAny) !== JSON.stringify(b.fieldAny)) return false;
	return true;
};

