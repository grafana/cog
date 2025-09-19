export const ConstantRefString = "AString";

export interface MyStruct {
	aString: "AString";
	optString?: "AString";
}

export const defaultMyStruct = (): MyStruct => ({
	aString: ConstantRefString,
	optString: ConstantRefString,
});
