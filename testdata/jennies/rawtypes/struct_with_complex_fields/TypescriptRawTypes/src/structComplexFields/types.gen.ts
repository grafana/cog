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

export const ConnectionPath = "straight";

export interface SomeOtherStruct {
	FieldAny: any;
}

export const defaultSomeOtherStruct = (): SomeOtherStruct => ({
	FieldAny: {},
});

