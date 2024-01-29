import * as externalPkg from '../externalPkg';


export interface Intersections extends SomeStruct, externalPkg.AnotherStruct {
	fieldString: string;
	fieldInteger: number;
}

export const defaultIntersections = (): Intersections => ({
	fieldString: "hello",
	fieldInteger: 32,
});

export interface SomeStruct {
	fieldBool: boolean;
}

export const defaultSomeStruct = (): SomeStruct => ({
	fieldBool: true,
});

