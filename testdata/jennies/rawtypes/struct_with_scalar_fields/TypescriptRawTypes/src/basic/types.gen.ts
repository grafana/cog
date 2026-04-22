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

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if (JSON.stringify(a.fieldAny) !== JSON.stringify(b.fieldAny)) return false;
	if (a.fieldBool !== b.fieldBool) return false;
	if (a.fieldBytes !== b.fieldBytes) return false;
	if (a.fieldString !== b.fieldString) return false;
	if (a.fieldStringWithConstantValue !== b.fieldStringWithConstantValue) return false;
	if (a.fieldFloat32 !== b.fieldFloat32) return false;
	if (a.fieldFloat64 !== b.fieldFloat64) return false;
	if (a.fieldUint8 !== b.fieldUint8) return false;
	if (a.fieldUint16 !== b.fieldUint16) return false;
	if (a.fieldUint32 !== b.fieldUint32) return false;
	if (a.fieldUint64 !== b.fieldUint64) return false;
	if (a.fieldInt8 !== b.fieldInt8) return false;
	if (a.fieldInt16 !== b.fieldInt16) return false;
	if (a.fieldInt32 !== b.fieldInt32) return false;
	if (a.fieldInt64 !== b.fieldInt64) return false;
	return true;
};

