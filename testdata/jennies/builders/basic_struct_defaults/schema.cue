package basic_struct_defaults

SomeStruct: {
	id: int64 | *42
	uid: string | *"default-uid"
	tags: [...string] | *["generated", "cog"]
	liveNow: bool | *true
}
