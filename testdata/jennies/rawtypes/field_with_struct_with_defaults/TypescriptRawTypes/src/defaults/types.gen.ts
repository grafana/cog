export interface NestedStruct {
	stringVal: string;
	intVal: number;
}

export const defaultNestedStruct = (): NestedStruct => ({
	stringVal: "",
	intVal: 0,
});

// equalsNestedStruct tests the equality of two `NestedStruct` objects.
export const equalsNestedStruct = (a: NestedStruct, b: NestedStruct): boolean => {
	if (a.stringVal !== b.stringVal) return false;
	if (a.intVal !== b.intVal) return false;
	return true;
};

export interface Struct {
	allFields: NestedStruct;
	partialFields: NestedStruct;
	emptyFields: NestedStruct;
	complexField: {
		uid: string;
		nested: {
			nestedVal: string;
		};
		array: string[];
	};
	partialComplexField: {
		uid: string;
		intVal: number;
	};
}

export const defaultStruct = (): Struct => ({
	allFields: { stringVal: "hello", intVal: 3, },
	partialFields: { stringVal: "", intVal: 3, },
	emptyFields: defaultNestedStruct(),
	complexField: { uid: "myUID", nested: { nestedVal: "nested", }, array: [
"hello",
], },
	partialComplexField: { uid: "", intVal: 0, },
});

// equalsStruct tests the equality of two `Struct` objects.
export const equalsStruct = (a: Struct, b: Struct): boolean => {
	if (!equalsNestedStruct(a.allFields, b.allFields)) return false;
	if (!equalsNestedStruct(a.partialFields, b.partialFields)) return false;
	if (!equalsNestedStruct(a.emptyFields, b.emptyFields)) return false;
	if (!equalsUnknown(a.complexField, b.complexField)) return false;
	if (!equalsUnknown(a.partialComplexField, b.partialComplexField)) return false;
	return true;
};

