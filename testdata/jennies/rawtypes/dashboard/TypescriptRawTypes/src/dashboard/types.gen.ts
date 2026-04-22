import * as cog from '../cog';


export interface Dashboard {
	title: string;
	panels?: Panel[];
}

export const defaultDashboard = (): Dashboard => ({
	title: "",
});

// equalsDashboard tests the equality of two `Dashboard` objects.
export const equalsDashboard = (a: Dashboard, b: Dashboard): boolean => {
	if (a.title !== b.title) return false;
	if ((a.panels === undefined) !== (b.panels === undefined)) return false;
	if (a.panels !== undefined) {
		if (a.panels.length !== b.panels!.length) return false;
		for (let i2 = 0; i2 < a.panels.length; i2++) {
			if (!equalsPanel(a.panels[i2], b.panels![i2])) return false;
		}
	}
	return true;
};

export interface DataSourceRef {
	type?: string;
	uid?: string;
}

export const defaultDataSourceRef = (): DataSourceRef => ({
});

// equalsDataSourceRef tests the equality of two `DataSourceRef` objects.
export const equalsDataSourceRef = (a: DataSourceRef, b: DataSourceRef): boolean => {
	if ((a.type === undefined) !== (b.type === undefined)) return false;
	if (a.type !== undefined) {
		if (a.type !== b.type!) return false;
	}
	if ((a.uid === undefined) !== (b.uid === undefined)) return false;
	if (a.uid !== undefined) {
		if (a.uid !== b.uid!) return false;
	}
	return true;
};

export interface FieldConfigSource {
	defaults?: FieldConfig;
}

export const defaultFieldConfigSource = (): FieldConfigSource => ({
});

// equalsFieldConfigSource tests the equality of two `FieldConfigSource` objects.
export const equalsFieldConfigSource = (a: FieldConfigSource, b: FieldConfigSource): boolean => {
	if ((a.defaults === undefined) !== (b.defaults === undefined)) return false;
	if (a.defaults !== undefined) {
		if (!equalsFieldConfig(a.defaults, b.defaults!)) return false;
	}
	return true;
};

export interface FieldConfig {
	unit?: string;
	custom?: any;
}

export const defaultFieldConfig = (): FieldConfig => ({
});

// equalsFieldConfig tests the equality of two `FieldConfig` objects.
export const equalsFieldConfig = (a: FieldConfig, b: FieldConfig): boolean => {
	if ((a.unit === undefined) !== (b.unit === undefined)) return false;
	if (a.unit !== undefined) {
		if (a.unit !== b.unit!) return false;
	}
	if ((a.custom === undefined) !== (b.custom === undefined)) return false;
	if (a.custom !== undefined) {
		if (JSON.stringify(a.custom) !== JSON.stringify(b.custom!)) return false;
	}
	return true;
};

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

// equalsPanel tests the equality of two `Panel` objects.
export const equalsPanel = (a: Panel, b: Panel): boolean => {
	if (a.title !== b.title) return false;
	if (a.type !== b.type) return false;
	if ((a.datasource === undefined) !== (b.datasource === undefined)) return false;
	if (a.datasource !== undefined) {
		if (!equalsDataSourceRef(a.datasource, b.datasource!)) return false;
	}
	if ((a.options === undefined) !== (b.options === undefined)) return false;
	if (a.options !== undefined) {
		if (JSON.stringify(a.options) !== JSON.stringify(b.options!)) return false;
	}
	if ((a.targets === undefined) !== (b.targets === undefined)) return false;
	if (a.targets !== undefined) {
		if (a.targets.length !== b.targets!.length) return false;
		for (let i2 = 0; i2 < a.targets.length; i2++) {
			if (a.targets[i2] !== b.targets![i2]) return false;
		}
	}
	if ((a.fieldConfig === undefined) !== (b.fieldConfig === undefined)) return false;
	if (a.fieldConfig !== undefined) {
		if (!equalsFieldConfigSource(a.fieldConfig, b.fieldConfig!)) return false;
	}
	return true;
};

