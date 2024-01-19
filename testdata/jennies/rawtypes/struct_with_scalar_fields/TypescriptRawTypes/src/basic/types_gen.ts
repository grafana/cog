// This
// is
// a
// comment
export interface SomeStruct {
	// Anything can go in there.
	// Really, anything.
	FieldAny: any;
	FieldBool: boolean;
	FieldBytes: string;
	FieldString: string;
	FieldStringWithConstantValue: "auto";
	FieldFloat32: number;
	FieldFloat64: number;
	FieldUint8: number;
	FieldUint16: number;
	FieldUint32: number;
	FieldUint64: number;
	FieldInt8: number;
	FieldInt16: number;
	FieldInt32: number;
	FieldInt64: number;
}

export const defaultSomeStruct = (): SomeStruct => ({
	FieldAny: {},
	FieldBool: false,
	FieldBytes: "",
	FieldString: "",
	FieldStringWithConstantValue: "auto",
	FieldFloat32: 0,
	FieldFloat64: 0,
	FieldUint8: 0,
	FieldUint16: 0,
	FieldUint32: 0,
	FieldUint64: 0,
	FieldInt8: 0,
	FieldInt16: 0,
	FieldInt32: 0,
	FieldInt64: 0,
});

