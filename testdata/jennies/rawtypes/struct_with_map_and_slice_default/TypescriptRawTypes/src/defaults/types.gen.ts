export interface SomeStruct {
	options?: Record<string, any>;
	items?: string[];
	extra: any;
}

export const defaultSomeStruct = (): SomeStruct => ({
	options: {
},
	items: [
],
	extra: {
},
});

// equalsSomeStruct tests the equality of two `SomeStruct` objects.
export const equalsSomeStruct = (a: SomeStruct, b: SomeStruct): boolean => {
	if ((a.options === undefined) !== (b.options === undefined)) return false;
	if (a.options !== undefined) {
		if (Object.keys(a.options).length !== Object.keys(b.options!).length) return false;
		for (const key2 in a.options) {
			if (JSON.stringify(a.options[key2]) !== JSON.stringify(b.options![key2])) return false;
		}
	}
	if ((a.items === undefined) !== (b.items === undefined)) return false;
	if (a.items !== undefined) {
		if (a.items.length !== b.items!.length) return false;
		for (let i2 = 0; i2 < a.items.length; i2++) {
			if (a.items[i2] !== b.items![i2]) return false;
		}
	}
	if (JSON.stringify(a.extra) !== JSON.stringify(b.extra)) return false;
	return true;
};

