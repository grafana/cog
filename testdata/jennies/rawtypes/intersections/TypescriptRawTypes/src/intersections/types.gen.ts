import * as externalPkg from '../externalPkg';


export interface Intersections extends SomeStruct, externalPkg.AnotherStruct {
	fieldString: string;
	fieldInteger: number;
}

export const defaultIntersections = (): Intersections => ({
	fieldString: "hello",
	fieldInteger: 32,
});

export interface SomeStruct {
	fieldBool: boolean;
}

export const defaultSomeStruct = (): SomeStruct => ({
	fieldBool: true,
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (a.fieldBool !== b.fieldBool) return false;
	return true;
};

// Base properties for all metrics
export interface Common {
	// The metric name
	name: string;
	// The metric type
	type: "counter" | "gauge";
	// The type of data the metric contains
	contains: "default" | "time";
}

export const defaultCommon = (): Common => ({
	name: "",
	type: "counter",
	contains: "default",
});

// equalsCommon tests the equality of two `Common` objects.
export const equalsCommon = (a: Common, b: Common): boolean => {
	if (a.name !== b.name) return false;
	if (a.type !== b.type) return false;
	if (a.contains !== b.contains) return false;
	return true;
};

// Counter metric combining common properties with specific values
export interface Counter extends Common {
	type: "counter";
	// Counter metric values
values: {
	// Total count of events
	count: number;
};
}

export const defaultCounter = (): Counter => ({
	type: "counter",
	values: {
	count: 0,
},
});

