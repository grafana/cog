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

export interface SomeOtherStruct {
	FieldAny: any;
}

export const defaultSomeOtherStruct = (): SomeOtherStruct => ({
	FieldAny: {},
});

