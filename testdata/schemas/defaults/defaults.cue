package defaults

VariableOption: {
  selected?: bool | string
  text: string | [...string]
	value: string | [...string]
}

TextVariable: {
  name: string | *""
  current: VariableOption | *{
    text: ""
    value: ["val"]
    selected: "maybe"
  }
  skipUrlSync: bool | *false
}
