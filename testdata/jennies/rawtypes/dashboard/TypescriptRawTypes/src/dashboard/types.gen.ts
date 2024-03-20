import * as cog from '../cog';


export interface Dashboard {
	title: string;
	panels?: Panel[];
}

export const defaultDashboard = (): Dashboard => ({
	title: "",
});

export interface DataSourceRef {
	type?: string;
	uid?: string;
}

export const defaultDataSourceRef = (): DataSourceRef => ({
});

export interface FieldConfigSource {
	defaults?: FieldConfig;
}

export const defaultFieldConfigSource = (): FieldConfigSource => ({
});

export interface FieldConfig {
	unit?: string;
	custom?: any;
}

export const defaultFieldConfig = (): FieldConfig => ({
});

export interface Panel {
	title: string;
	type: string;
	datasource?: DataSourceRef;
	options?: any;
	targets?: cog.Dataquery[];
	fieldConfig?: FieldConfigSource;
}

export const defaultPanel = (): Panel => ({
	title: "",
	type: "",
});

