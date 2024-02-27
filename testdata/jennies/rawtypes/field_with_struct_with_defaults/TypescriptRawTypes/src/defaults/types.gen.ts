export interface NestedStruct {
	stringVal: string;
	intVal: number;
}

export const defaultNestedStruct = (): NestedStruct => ({
	stringVal: "",
	intVal: 0,
});

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

