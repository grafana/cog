export interface Query {
	expr: string;
	instant?: boolean;
	_implementsDataqueryVariant(): void;
}

export const defaultQuery = (): Query => ({
	expr: "",
	_implementsDataqueryVariant: () => {},
});

