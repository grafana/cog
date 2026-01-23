package equality

#Direction: "top" | "bottom" | "left" | "right"

#Variable: {
    name: string
}

container: {
	stringField: string
	intField: int
	enumField: #Direction
	refField: #Variable
}

optionals: {
	stringField?: string
	enumField?: #Direction
	refField?: #Variable
	byteField?: bytes
}

arrays: {
    ints: [...int]
    strings: [...string]
    arrayOfArray: [...[...string]]
    refs: [...#Variable]
    anonymousStructs: [...{
    	inner: string
    }]
    arrayOfAny: [...]
}

maps: {
    ints: [string]: int
    strings: [string]: string
    refs: [string]: #Variable
    anonymousStructs: [string]: {
    	inner: string
    }
    stringToAny: {
				[string]: _
    }
}
