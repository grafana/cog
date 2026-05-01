export interface Options {
	timeseries_option: string;
}

export const defaultOptions = (): Options => ({
	timeseries_option: "",
});

// equalsOptions tests the equality of two `Options` objects.
export const equalsOptions = (a: Options, b: Options): boolean => {
	if (a.timeseriesOption !== b.timeseriesOption) return false;
	return true;
};

export interface FieldConfig {
	timeseries_field_config_option: string;
}

export const defaultFieldConfig = (): FieldConfig => ({
	timeseries_field_config_option: "",
});

// equalsFieldConfig tests the equality of two `FieldConfig` objects.
export const equalsFieldConfig = (a: FieldConfig, b: FieldConfig): boolean => {
	if (a.timeseriesFieldConfigOption !== b.timeseriesFieldConfigOption) return false;
	return true;
};

