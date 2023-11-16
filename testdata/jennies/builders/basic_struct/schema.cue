package basic_struct

// SomeStruct, to hold data.
SomeStruct: {
	// id identifies something. Weird, right?
	id: int64
	uid: string
	tags: [...string]
	// This thing could be live.
	// Or maybe not.
	liveNow: bool
}
