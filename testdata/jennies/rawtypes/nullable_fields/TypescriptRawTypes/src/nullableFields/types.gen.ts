export interface Struct {
	a: MyObject | null;
	b?: MyObject | null;
	c: string | null;
	d: string[] | null;
	e: Record<string, string | null>;
	f: {
		a: string;
	} | null;
	g: "hey" | null;
}

export const defaultStruct = (): Struct => ({
	a: defaultMyObject(),
	c: "",
	d: [],
	e: {},
	f: {
	a: "",
},
	g: ConstantRef,
});

// equalsStruct tests the equality of two `Struct` objects.
export const equalsStruct = (a: Struct, b: Struct): boolean => {
	if (a.a !== b.a) return false;
	if ((a.b === undefined) !== (b.b === undefined)) return false;
	if (a.b !== undefined) {
		if (a.b !== b.b!) return false;
	}
	if (a.c !== b.c) return false;
	if (a.d !== b.d) return false;
	if (Object.keys(a.e).length !== Object.keys(b.e).length) return false;
	for (const key1 in a.e) {
		if (a.e[key1] !== b.e[key1]) return false;
	}
	if (a.f !== b.f) return false;
	if (a.g !== b.g) return false;
	return true;
};

export interface MyObject {
	field: string;
}

export const defaultMyObject = (): MyObject => ({
	field: "",
});

// equalsMyObject tests the equality of two `MyObject` objects.
export const equalsMyObject = (a: MyObject, b: MyObject): boolean => {
	if (a.field !== b.field) return false;
	return true;
};

export const ConstantRef = "hey";

