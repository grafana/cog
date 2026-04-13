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

export interface MyObject {
	field: string;
}

export const defaultMyObject = (): MyObject => ({
	field: "",
});

export const ConstantRef = "hey";

