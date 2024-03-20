export interface Options {
	timeseries_option: string;
}

export const defaultOptions = (): Options => ({
	timeseries_option: "",
});

export interface FieldConfig {
	timeseries_field_config_option: string;
}

export const defaultFieldConfig = (): FieldConfig => ({
	timeseries_field_config_option: "",
});

