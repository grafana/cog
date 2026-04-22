export const ConstantRefString = "AString";

export interface MyStruct {
	aString: "AString";
	optString?: "AString";
}

export const defaultMyStruct = (): MyStruct => ({
	aString: ConstantRefString,
	optString: ConstantRefString,
});

// equalsMyStruct tests the equality of two `MyStruct` objects.
export const equalsMyStruct = (a: MyStruct, b: MyStruct): boolean => {
	if (a.aString !== b.aString) return false;
	if ((a.optString === undefined) !== (b.optString === undefined)) return false;
	if (a.optString !== undefined) {
		if (a.optString !== b.optString!) return false;
	}
	return true;
};

