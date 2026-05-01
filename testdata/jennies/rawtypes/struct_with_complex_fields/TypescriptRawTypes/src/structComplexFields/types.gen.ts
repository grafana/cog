// This struct does things.
export interface SomeStruct {
	FieldRef: SomeOtherStruct;
	FieldDisjunctionOfScalars: string | boolean;
	FieldMixedDisjunction: string | SomeOtherStruct;
	FieldDisjunctionWithNull: string | null;
	Operator: ">" | "<";
	FieldArrayOfStrings: string[];
	FieldMapOfStringToString: Record<string, string>;
	FieldAnonymousStruct: {
		FieldAny: any;
	};
	fieldRefToConstant: "straight";
}

export const defaultSomeStruct = (): SomeStruct => ({
	FieldRef: defaultSomeOtherStruct(),
	FieldDisjunctionOfScalars: "",
	FieldMixedDisjunction: "",
	FieldDisjunctionWithNull: "",
	Operator: ">",
	FieldArrayOfStrings: [],
	FieldMapOfStringToString: {},
	FieldAnonymousStruct: {
	FieldAny: {},
},
	fieldRefToConstant: ConnectionPath,
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (!equalsSomeOtherStruct(a.fieldRef, b.fieldRef)) return false;
	if (a.fieldDisjunctionOfScalars !== b.fieldDisjunctionOfScalars) return false;
	if (a.fieldMixedDisjunction !== b.fieldMixedDisjunction) return false;
	if (a.fieldDisjunctionWithNull !== b.fieldDisjunctionWithNull) return false;
	if (a.operator !== b.operator) return false;
	if (a.fieldArrayOfStrings.length !== b.fieldArrayOfStrings.length) return false;
	for (let i1 = 0; i1 < a.fieldArrayOfStrings.length; i1++) {
		if (a.fieldArrayOfStrings[i1] !== b.fieldArrayOfStrings[i1]) return false;
	}
	if (Object.keys(a.fieldMapOfStringToString).length !== Object.keys(b.fieldMapOfStringToString).length) return false;
	for (const key1 in a.fieldMapOfStringToString) {
		if (a.fieldMapOfStringToString[key1] !== b.fieldMapOfStringToString[key1]) return false;
	}
	if (!equalsUnknown(a.fieldAnonymousStruct, b.fieldAnonymousStruct)) return false;
	if (a.fieldRefToConstant !== b.fieldRefToConstant) return false;
	return true;
};

export const ConnectionPath = "straight";

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

